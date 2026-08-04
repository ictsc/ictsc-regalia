package domain

import "context"

// Txは処理をトランザクション内で実行するためのインターフェースを表す
type Tx[E any] interface {
	RunInTx(ctx context.Context, f func(E) error) error
}
