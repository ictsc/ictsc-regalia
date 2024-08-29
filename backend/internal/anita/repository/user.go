package repository

import (
	"context"

	"github.com/ictsc/ictsc-outlands/backend/internal/anita/domain"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/domain/value"
)

// UserRepository ユーザーリポジトリ
type UserRepository interface {
	SelectUser(ctx context.Context, id value.UserID) (*domain.User, error)
	SelectUsers(ctx context.Context) ([]*domain.User, error)
	SelectUsersByTeamID(ctx context.Context, teamID value.TeamID) ([]*domain.User, error)
	UpsertUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id value.UserID) error
}
