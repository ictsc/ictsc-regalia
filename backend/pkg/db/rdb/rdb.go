package rdb

import (
	"context"
	"database/sql"
)

// Tx RDBトランザクションインターフェース
type Tx interface {
	Do(ctx context.Context, options *sql.TxOptions, callBack func(context.Context) error) error
}
