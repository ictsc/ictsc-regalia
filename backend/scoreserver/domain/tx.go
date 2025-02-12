package domain

import (
	"context"

	"github.com/cockroachdb/errors"
)

type Tx[E any] interface {
	RunInTx(ctx context.Context, f func(E) error) error
}

type TxFunc[E any] func(context.Context, func(E) error) error

func (t TxFunc[E]) RunInTx(ctx context.Context, f func(E) error) error {
	return t(ctx, f)
}

func RunTx[V any, E any](ctx context.Context, tx Tx[E], fun func(E) (V, error)) (V, error) {
	var val V
	if tx == nil {
		return val, errors.New("tx is nil")
	}
	err := tx.RunInTx(ctx, func(e E) error {
		var err error
		val, err = fun(e)
		return err
	})
	return val, err
}
