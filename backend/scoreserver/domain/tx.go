package domain

import "context"

type TxFunc[E any] func(context.Context, func(eff E) error) error
