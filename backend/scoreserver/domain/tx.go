package domain

import "context"

type Tx[E any] interface {
	RunInTx(ctx context.Context, f func(E) error) error
}

type TxFunc[E any] func(context.Context, func(E) error) error

func (t TxFunc[E]) RunInTx(ctx context.Context, f func(E) error) error {
	return t(ctx, f)
}
