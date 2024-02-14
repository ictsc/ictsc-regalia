package bun

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
)

type txsKey struct{}

var key = txsKey{}

// Do is a func to execute db operations in transaction
func (repo *DB) Do(ctx context.Context, options *sql.TxOptions, callBack func(context.Context) error) error {
	var (
		idb bun.IDB
		ok  bool
	)

	idb, ok = ctx.Value(key).(bun.Tx)
	if !ok {
		idb = repo.DB
	}

	tx, err := idb.BeginTx(ctx, options)
	if err != nil {
		return err
	}

	defer tx.Rollback()
	ctx = context.WithValue(ctx, key, tx)

	err = callBack(ctx)
	if err != nil {
		return err
	}

	err = ctx.Err()
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (db *DB) getIDB(ctx context.Context) bun.IDB {
	tx, ok := ctx.Value(key).(bun.Tx)
	if !ok {
		return db.DB
	}

	return tx
}
