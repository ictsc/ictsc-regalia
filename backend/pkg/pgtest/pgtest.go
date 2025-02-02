package pgtest

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"github.com/cockroachdb/errors"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	postgresImage = "postgres:17"
	schemaFile    = "schema.sql"
)

// SetupDB はテスト用の DB を用意します。
//
// txdb を用いているため高速ですが，本物の DB とは異なる挙動をする可能性があります。
// 問題が起きた場合は SetupTrueDB を用いてください
func SetupDB(tb testing.TB) *sqlx.DB {
	tb.Helper()

	ctx := context.Background()

	ctr := getContainer(ctx)

	connString, err := ctr.ConnectionString(ctx)
	if err != nil {
		tb.Fatalf("Failed to get connection string: %v", err)
		return nil
	}

	db := sqlx.NewDb(sql.OpenDB(txdb.New("pgx", connString)), "pgx")
	tb.Cleanup(func() {
		if err := db.Close(); err != nil {
			log.Printf("Failed to close DB: %v\n", err)
		}
	})
	return db
}

func SetupTrueDB(tb testing.TB) *sqlx.DB {
	tb.Helper()

	// TODO: startContainer を使って DB を起動し，tb.Cleanup で停止する
	tb.Fatal("unimplemented")
	return nil
}

var (
	container *postgres.PostgresContainer
	startOnce sync.Once
)

func getContainer(ctx context.Context) *postgres.PostgresContainer {
	startOnce.Do(func() {
		var err error
		container, err = startContainer(ctx)
		if err != nil {
			panic(err)
		}
	})
	return container
}

func startContainer(ctx context.Context) (*postgres.PostgresContainer, error) {
	execDir, err := os.Getwd()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	baseDir := execDir
	for {
		goModPath := filepath.Join(baseDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			break
		}

		parentDir := filepath.Dir(baseDir)
		if parentDir == baseDir {
			break
		}

		baseDir = parentDir
	}

	container, err := postgres.Run(ctx,
		postgresImage,
		postgres.WithInitScripts(filepath.Join(baseDir, schemaFile)),
		testcontainers.WithWaitStrategy(
			//nolint:mnd
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		))
	if err != nil {
		return nil, errors.Wrap(err, "failed to start Postgres container")
	}
	return container, nil
}
