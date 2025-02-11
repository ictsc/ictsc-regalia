package pg

import (
	"context"
	"database/sql"
	"iter"

	"github.com/cockroachdb/errors"
	"github.com/gofrs/uuid/v5"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
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

var _ domain.UserCreator = (*RepositoryTx)(nil)

// CreateUser - ユーザ+プロフィールを作成する
//
// 複数のテーブルに関与するのでトランザクション内でしか呼び出せない
func (r *RepositoryTx) CreateUser(ctx context.Context, profile *domain.UserProfileData) error {
	if _, err := sqlx.NamedExecContext(ctx, r.ext, `
		INSERT INTO users (id, name, created_at) VALUES (:id, :name, NOW())
	`, (*userRow)(profile.User)); err != nil {
		if pgErr := new(pgconn.PgError); errors.As(err, &pgErr) {
			// 一意制約違反
			if pgErr.Code == "23505" {
				return domain.NewAlreadyExistsError("user", nil)
			}
		}
		return errors.Wrap(err, "failed to insert into users")
	}

	if _, err := sqlx.NamedExecContext(ctx, r.ext, `
		INSERT INTO user_profiles (user_id, display_name, created_at, updated_at)
		VALUES (:user_id, :display_name, NOW(), NOW())`,
		newUserProfileRow(profile),
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
			if pgErr.Code == "23505" {
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
	userProfileRow struct {
		UserID      uuid.UUID `db:"user_id"`
		DisplayName string    `db:"display_name"`
	}
)

func newUserProfileRow(profile *domain.UserProfileData) *userProfileRow {
	return &userProfileRow{
		UserID:      profile.User.ID,
		DisplayName: profile.DisplayName,
	}
}
