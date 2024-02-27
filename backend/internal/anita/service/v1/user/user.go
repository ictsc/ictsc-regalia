// Package user ユーザーサービス実装
package user

import (
	"github.com/ictsc/ictsc-outlands/backend/internal/anita/repository"
	srv "github.com/ictsc/ictsc-outlands/backend/internal/anita/service"
)

type service struct {
	repo repository.UserRepository
}

var _ srv.UserService = (*service)(nil)

func NewService(repo repository.UserRepository) *service {
	return &service{repo: repo}
}
