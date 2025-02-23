package pg

import (
	"context"
	"database/sql"
	"iter"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

var _ domain.UserLister = (*repo)(nil)

func (r *repo) ListUsers(ctx context.Context, filter domain.UserListFilter) iter.Seq2[*domain.UserData, error] {
	return func(yield func(*domain.UserData, error) bool) {
		cond := "TRUE"
		var args []any
		if filter.Name != "" {
			cond += " AND name = ?"
			args = append(args, filter.Name)
		}
		rows, err := r.ext.QueryxContext(ctx, r.ext.Rebind(`
			SELECT id, name FROM users
			WHERE `+cond,
		), args...)
		if err != nil {
			yield(nil, err)
			return
		}
		defer func() { _ = rows.Close() }()

		for rows.Next() {
			var row userRow
			if err := rows.StructScan(&row); err != nil {
				yield(nil, err)
				return
			}
			if !yield((*domain.UserData)(&row), nil) {
				return
			}
		}
	}
}

var _ domain.DiscordLinkedUserGetter = (*repo)(nil)

func (r *repo) GetDiscordLinkedUser(ctx context.Context, discordUserID int64) (*domain.UserData, error) {
	var row userRow
	if err := sqlx.GetContext(ctx, r.ext, &row, `
		SELECT u.id, u.name
		FROM users AS u
		JOIN discord_users AS d ON d.user_id = u.id
		WHERE d.discord_user_id = $1
	`, discordUserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("discord linked user", nil)
		}
		return nil, errors.Wrap(err, "failed to get discord linked user")
	}

	return (*domain.UserData)(&row), nil
}

var _ domain.UserProfileReader = (*repo)(nil)

func (r *repo) GetUserProfileByID(ctx context.Context, userID uuid.UUID) (*domain.UserProfileData, error) {
	var row struct {
		User    userRow    `db:"u"`
		Profile profileRow `db:"p"`
	}
	if err := sqlx.GetContext(ctx, r.ext, &row, `
		SELECT `+userColumns.As("u")+`, `+profileColumns.As("p")+`
		FROM users AS u
		INNER JOIN user_profiles AS p ON u.id = p.user_id
		WHERE u.id = $1`, userID,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("user profile", nil)
		}
		return nil, errors.Wrap(err, "failed to get user profile")
	}
	return &domain.UserProfileData{
		User:    (*domain.UserData)(&row.User),
		Profile: (*domain.ProfileData)(&row.Profile),
	}, nil
}

var _ domain.TeamMemberProfileReader = (*repo)(nil)

var (
	teamMemberProfileSelectQuery = `
SELECT
	` + userColumns.As("u") + `,
	` + profileColumns.As("p") + `,
	` + teamColumns.As("t") + `,
	d.discord_user_id
FROM users AS u
INNER JOIN user_profiles AS p ON u.id = p.user_id
INNER JOIN team_members AS tm ON u.id = tm.user_id
INNER JOIN teams AS t ON tm.team_id = t.id
LEFT JOIN discord_users AS d ON u.id = d.user_id`
	teamMemberProfileSelectQueryByTeamID = teamMemberProfileSelectQuery + "\nWHERE t.id = ?"
)

func (r *repo) ListTeamMembers(ctx context.Context) ([]*domain.TeamMemberProfileData, error) {
	return r.listTeamMembers(ctx, teamMemberProfileSelectQuery)
}

func (r *repo) ListTeamMembersByTeamID(ctx context.Context, teamID uuid.UUID) ([]*domain.TeamMemberProfileData, error) {
	return r.listTeamMembers(ctx, teamMemberProfileSelectQueryByTeamID, teamID)
}

func (r *repo) listTeamMembers(ctx context.Context, query string, args ...any) ([]*domain.TeamMemberProfileData, error) {
	rows, err := r.ext.QueryxContext(ctx, r.ext.Rebind(query), args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list team members")
	}
	defer func() { _ = rows.Close() }()

	var members []*domain.TeamMemberProfileData

	for rows.Next() {
		var row struct {
			User          userRow       `db:"u"`
			Profile       profileRow    `db:"p"`
			Team          teamRow       `db:"t"`
			DiscordUserID sql.NullInt64 `db:"discord_user_id"`
		}
		if err := rows.StructScan(&row); err != nil {
			return nil, errors.Wrap(err, "failed to scan team member profile")
		}
		members = append(members, &domain.TeamMemberProfileData{
			User:          (*domain.UserData)(&row.User),
			Profile:       (*domain.ProfileData)(&row.Profile),
			Team:          (*domain.TeamData)(&row.Team),
			DiscordUserID: row.DiscordUserID.Int64,
		})
	}
	return members, nil
}

var _ domain.UserCreator = (*RepositoryTx)(nil)

// CreateUser - ユーザ+プロフィールを作成する
//
// 複数のテーブルに関与するのでトランザクション内でしか呼び出せない
func (r *RepositoryTx) CreateUser(ctx context.Context, profile *domain.UserProfileData) error {
	if _, err := sqlx.NamedExecContext(ctx, r.ext, `
		INSERT INTO users (id, name, created_at) VALUES (:id, :name, NOW())
	`, (*userRow)(profile.User)); err != nil {
		// 一意制約違反
		if pgErr := new(pgconn.PgError); errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if pgErr.ConstraintName == "users_name_key" {
				return errors.WithStack(domain.ErrDuplicateUserName)
			}
			return domain.NewAlreadyExistsError("user", errors.Newf("constraint: %s", pgErr.ConstraintName))
		}
		return errors.Wrap(err, "failed to insert into users")
	}

	if _, err := sqlx.NamedExecContext(ctx, r.ext, `
		INSERT INTO user_profiles (user_id, display_name, created_at, updated_at)
		VALUES (:user_id, :display_name, NOW(), NOW())`,
		struct {
			UserID uuid.UUID `db:"user_id"`
			profileRow
		}{profile.User.ID, (profileRow)(*profile.Profile)},
	); err != nil {
		return errors.Wrap(err, "failed to insert into user_profiles")
	}

	return nil
}

var _ domain.DiscordUserLinker = (*RepositoryTx)(nil)

func (r *RepositoryTx) LinkDiscordUser(ctx context.Context, userID uuid.UUID, discordUserID int64) error {
	if _, err := r.ext.ExecContext(ctx, `
		INSERT INTO discord_users (user_id, discord_user_id, linked_at)
		VALUES ($1, $2, NOW())
	`, userID, discordUserID); err != nil {
		if pgErr := new(pgconn.PgError); errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return domain.NewAlreadyExistsError("discord user", nil)
			}
		}
		return domain.WrapAsInternal(err, "failed to insert into discord_user")
	}
	return nil
}

type (
	userRow struct {
		ID   uuid.UUID `db:"id"`
		Name string    `db:"name"`
	}
	profileRow struct {
		DisplayName string `db:"display_name"`
	}
)

var (
	userColumns    = columns([]string{"id", "name"})
	profileColumns = columns([]string{"display_name"})
)
