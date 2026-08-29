package memory

import (
	"sync"

	"github.com/ictsc/ictsc-regalia/backend/internal/domain"
)

// Repositoryはデータをメモリ上に保持するRepositoryのスタブ実装を表す
type Repository struct {
	mu    sync.RWMutex
	teams map[int64]*domain.Team
}

var _ domain.TeamRepository = (*Repository)(nil)

// NewRepositoryは空のRepositoryを作成する
func NewRepository() *Repository {
	return &Repository{
		teams: make(map[int64]*domain.Team),
	}
}
