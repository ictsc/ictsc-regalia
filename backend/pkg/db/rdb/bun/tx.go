package bun

import (
	"context"
	"database/sql"

	"github.com/ictsc/ictsc-outlands/backend/pkg/errors"
	"github.com/uptrace/bun"
)

type txKey bool

const key = txKey(false)

// Do トランザクションの中で処理を実行
func (d *DB) Do(ctx context.Context, options *sql.TxOptions, callBack func(context.Context) error) error {
	var (
		idb bun.IDB
		ok  bool
	)

	idb, ok = ctx.Value(key).(bun.Tx)
	if !ok {
		idb = d.db
	}

	tx, err := idb.BeginTx(ctx, options)
	if err != nil {
		return errors.Wrap(errors.ErrUnknown, err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	ctx = context.WithValue(ctx, key, tx)

	err = callBack(ctx)
	if err != nil {
		return err
	}

	err = ctx.Err()
	if err != nil {
		return errors.Wrap(errors.ErrUnknown, err)
	}

	return errors.Wrap(errors.ErrUnknown, tx.Commit())
}

// GetIDB bun.IDBインスタンスを取得
func (d *DB) GetIDB(ctx context.Context) bun.IDB { // nolint:ireturn
	tx, ok := ctx.Value(key).(bun.Tx)
	if !ok {
		return d.db
	}

	return tx
}
