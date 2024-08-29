package bun

import (
	"context"

	"github.com/ictsc/ictsc-outlands/backend/internal/anita/domain"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/domain/value"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/repository"
	"github.com/ictsc/ictsc-outlands/backend/pkg/db/rdb/bun"
	"github.com/ictsc/ictsc-outlands/backend/pkg/errors"
)

// UserRepository ユーザーリポジトリ
type UserRepository struct {
	db *bun.DB
}

var _ repository.UserRepository = (*UserRepository)(nil)

// NewUserRepository ユーザーリポジトリを生成する
func NewUserRepository(db *bun.DB) *UserRepository {
	return &UserRepository{db: db}
}

// SelectUser ユーザーを取得する
func (repo *UserRepository) SelectUser(ctx context.Context, id value.UserID) (*domain.User, error) {
	db := repo.db.GetIDB(ctx)
	user := new(User)

	exists, err := db.NewSelect().Model(user).Where("id = ?", id.String()).Exists(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.ErrUnknown, err)
	}

	if !exists {
		return nil, errors.Wrap(errors.ErrNotFound, nil)
	}

	err = db.NewSelect().Model(user).Where("id = ?", id.String()).Scan(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.ErrUnknown, err)
	}

	return convertToDomainUser(user)
}

// SelectUsers ユーザーを取得する
func (repo *UserRepository) SelectUsers(ctx context.Context) ([]*domain.User, error) {
	db := repo.db.GetIDB(ctx)
	users := make([]*User, 0)

	err := db.NewSelect().Model(&users).Scan(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.ErrUnknown, err)
	}

	return convertToDomainUsers(users)
}

// SelectUsersByTeamID チームIDからユーザーを取得する
func (repo *UserRepository) SelectUsersByTeamID(ctx context.Context, teamID value.TeamID) ([]*domain.User, error) {
	db := repo.db.GetIDB(ctx)
	users := make([]*User, 0)

	err := db.NewSelect().Model(&users).Where("team_id = ?", teamID.String()).Scan(ctx)
	if err != nil {
		return nil, errors.Wrap(errors.ErrUnknown, err)
	}

	return convertToDomainUsers(users)
}

// UpsertUser ユーザーを挿入または更新する
func (repo *UserRepository) UpsertUser(ctx context.Context, user *domain.User) error {
	db := repo.db.GetIDB(ctx)
	convUser := convertFromDomainUser(user)

	_, err := db.NewInsert().Model(convUser).On("DUPLICATE KEY UPDATE").Exec(ctx)
	if err != nil {
		return errors.Wrap(errors.ErrUnknown, err)
	}

	return nil
}

// DeleteUser ユーザーを削除する
func (repo *UserRepository) DeleteUser(ctx context.Context, id value.UserID) error {
	db := repo.db.GetIDB(ctx)

	_, err := db.NewDelete().Model((*User)(nil)).Where("id = ?", id.String()).Exec(ctx)
	if err != nil {
		return errors.Wrap(errors.ErrUnknown, err)
	}

	return nil
}
