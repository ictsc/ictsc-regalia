package pg

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
	repo
}

type RepositoryTx struct {
	repo
}

type repo struct {
	ext sqlx.ExtContext
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db:   db,
		repo: repo{ext: db},
	}
}

//nolint:varnamelen // API としての名前なので許容
func (r *Repository) RunTx(ctx context.Context, fn func(tx *RepositoryTx) error) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin tx")
	}

	if err := fn(&RepositoryTx{repo: repo{ext: tx}}); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = errors.Join(err, errors.Wrap(rerr, "failed to rollback tx"))
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit tx")
	}

	return nil
}

func Tx[E any](r *Repository, weaker func(*RepositoryTx) E) domain.Tx[E] {
	return domain.TxFunc[E](func(ctx context.Context, f func(E) error) error {
		return r.RunTx(ctx, func(tx *RepositoryTx) error {
			return f(weaker(tx))
		})
	})
}
