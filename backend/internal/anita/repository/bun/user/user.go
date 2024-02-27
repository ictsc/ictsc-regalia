// Package user ユーザーリポジトリ実装
package user

import (
	repo "github.com/ictsc/ictsc-outlands/backend/internal/anita/repository"
	"github.com/ictsc/ictsc-outlands/backend/pkg/db/rdb/bun"
)

type repository struct {
	db *bun.DB
}

var _ repo.UserRepository = (*repository)(nil)

func NewRepository(db *bun.DB) *repository {
	return &repository{db: db}
}
