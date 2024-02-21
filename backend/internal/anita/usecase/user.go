package usecase

import (
	"context"

	"github.com/ictsc/ictsc-outlands/backend/internal/domain"
	"github.com/oklog/ulid/v2"
)

// UserUsecase ユーザーユースケース
type UserUsecase interface {
	ReadUser(ctx context.Context, id ulid.ULID) (*domain.User, error)
	ReadUsers(ctx context.Context) ([]*domain.User, error)
	ReadUsersByTeamID(ctx context.Context, teamID ulid.ULID) ([]*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
	UpdateUser(ctx context.Context, id ulid.ULID, name string) (*domain.User, error)
	DeleteUser(ctx context.Context, id ulid.ULID) error
}
