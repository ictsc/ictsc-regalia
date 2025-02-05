package pg

import (
	"context"
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
				return domain.NewError(domain.ErrTypeAlreadyExists, errors.New("user already exists"))
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
