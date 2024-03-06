package service

import (
	"context"

	"github.com/ictsc/ictsc-outlands/backend/internal/anita/domain"
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/domain/value"
)

// UserService ユーザーサービス
type UserService interface {
	ReadUser(ctx context.Context, id value.UserID) (*domain.User, error)
	ReadUsers(ctx context.Context) ([]*domain.User, error)
	CreateUser(ctx context.Context, id value.UserID, name value.UserName, code value.TeamInvitationCode) (*domain.User, error)
	UpdateUser(ctx context.Context, id value.UserID, name value.UserName) (*domain.User, error)
	DeleteUser(ctx context.Context, id value.UserID) error
}
