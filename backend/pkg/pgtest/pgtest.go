package pgtest

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"
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
	postgresImage = "postgres:13"
	schemaFile    = "schema.sql"
)

var (
	postgresContainer *postgres.PostgresContainer
)

type RunFunc func() int

// WrapRun は関数をラップして Postgres のライフサイクルを管理します。
//
// TestMain で呼び出すことを想定しています。
func WrapRun(run RunFunc) RunFunc {
	return func() int {
		ctx := context.Background()

		ctr, err := startContainer(ctx)
		if err != nil {
			log.Printf("Failed to start Postgres container: %v\n", err)
			return 1
		}
		defer func() {
			if err := testcontainers.TerminateContainer(ctr); err != nil {
				log.Printf("Failed to terminate Postgres container: %v\n", err)
			}
		}()

		postgresContainer = ctr
		return run()
	}
}

// SetupDB はテスト用の DB を用意します。
//
// txdb を用いているため高速ですが，本物の DB とは異なる挙動をする可能性があります。
// 問題が起きた場合は SetupTrueDB を用いてください
func SetupDB(tb testing.TB) (*sqlx.DB, bool) {
	tb.Helper()

	ctx := context.Background()

	if postgresContainer == nil {
		tb.Errorf("Postgres Container is not initialized. May forgot to call WrapRun?")
		return nil, false
	}

	connString, err := postgresContainer.ConnectionString(ctx)
	if err != nil {
		tb.Errorf("Failed to get connection string: %v", err)
		return nil, false
	}

	return sqlx.NewDb(sql.OpenDB(txdb.New("pgx", connString)), "pgx"), true
}

func SetupTrueDB(tb testing.TB) (*sqlx.DB, bool) {
	tb.Helper()

	// TODO: startContainer を使って DB を起動し，tb.Cleanup で停止する
	tb.Error("unimplemented")
	return nil, false
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
	return container, errors.WithStack(err)
}
