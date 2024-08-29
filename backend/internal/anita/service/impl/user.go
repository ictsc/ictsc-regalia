package impl

import (
	"context"

	"github.com/ictsc/ictsc-outlands/backend/internal/anita/domain"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/domain/value"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/repository"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/service"
	"github.com/ictsc/ictsc-outlands/backend/pkg/db/rdb"
	"github.com/ictsc/ictsc-outlands/backend/pkg/errors"
)

// UserService ユーザーサービス
type UserService struct {
	tx    rdb.Tx
	repo  repository.UserRepository
	tRepo repository.TeamRepository
}

var _ service.UserService = (*UserService)(nil)

// NewUserService ユーザーサービスを生成する
func NewUserService(tx rdb.Tx, repo repository.UserRepository, tRepo repository.TeamRepository) *UserService {
	return &UserService{tx: tx, repo: repo, tRepo: tRepo}
}

func (s *UserService) exists(ctx context.Context, id value.UserID) bool {
	_, err := s.repo.SelectUser(ctx, id)

	return err == nil || !errors.IsFlag(err, errors.ErrNotFound)
}

// ReadUser ユーザーを取得する
func (s *UserService) ReadUser(ctx context.Context, id value.UserID) (*domain.User, error) {
	return s.repo.SelectUser(ctx, id)
}

// ReadUsers ユーザーを取得する
func (s *UserService) ReadUsers(ctx context.Context) ([]*domain.User, error) {
	return s.repo.SelectUsers(ctx)
}

// CreateUser ユーザーを作成する
func (s *UserService) CreateUser(ctx context.Context, id value.UserID, name value.UserName, code value.TeamInvitationCode) (*domain.User, error) {
	var user *domain.User

	err := s.tx.Do(ctx, nil, func(ctx context.Context) error {
		if s.exists(ctx, id) {
			return errors.New(errors.ErrAlreadyExists, "User already exists")
		}

		team, err := s.tRepo.SelectTeamByInvitationCode(ctx, code)
		if err != nil {
			return err
		}

		teamID := team.ID()
		user = domain.NewUser(id, name, teamID)

		err = s.repo.UpsertUser(ctx, user)
		if err != nil {
			return err
		}

		err = team.DecrementCodeRemaining()
		if err != nil {
			return err
		}

		return s.tRepo.UpsertTeam(ctx, team)
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser ユーザーを更新する
func (s *UserService) UpdateUser(ctx context.Context, id value.UserID, name value.UserName) (*domain.User, error) {
	var (
		user *domain.User
		err  error
	)

	err = s.tx.Do(ctx, nil, func(ctx context.Context) error {
		user, err = s.repo.SelectUser(ctx, id)
		if err != nil {
			return err
		}

		user.SetName(name)

		return s.repo.UpsertUser(ctx, user)
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser ユーザーを削除する
func (s *UserService) DeleteUser(ctx context.Context, id value.UserID) error {
	err := s.tx.Do(ctx, nil, func(ctx context.Context) error {
		if !s.exists(ctx, id) {
			return errors.New(errors.ErrNotFound, "User not found")
		}

		return s.repo.DeleteUser(ctx, id)
	})
	if err != nil {
		return err
	}

	return nil
}
