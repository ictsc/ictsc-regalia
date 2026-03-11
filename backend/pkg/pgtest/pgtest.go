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

// SetupDB はテスト用の DB を用意します。
//
// txdb を用いているため高速ですが，本物の DB とは異なる挙動をする可能性があります。
// 問題が起きた場合は SetupTrueDB を用いてください
func SetupDB(tb testing.TB) *sqlx.DB {
	tb.Helper()

	ctr, err := getContainer()
	if err != nil {
		tb.Fatalf("Failed to start Postgres container: %v", err)
		return nil
	}

	ctx := tb.Context()
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
	startErr  error
	startOnce sync.Once
)

func getContainer() (*postgres.PostgresContainer, error) {
	startOnce.Do(func() {
		for i := 0; i < containerStartRetries; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), containerStartTimeout)
			ctr, err := startContainer(ctx)
			cancel()
			if err == nil {
				container = ctr
				startErr = nil
				return
			}

			startErr = err
			if i < containerStartRetries-1 {
				time.Sleep(containerRetryInterval)
			}
		}
	})

	return container, startErr
}

const (
	postgresImage          = "postgres:17"
	schemaFile             = "schema.sql"
	viewFile               = "view.sql"
	seedFile               = "seed.sql"
	containerStartTimeout  = 60 * time.Second
	containerStartRetries  = 3
	containerRetryInterval = 2 * time.Second
)

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
		postgres.WithInitScripts(
			filepath.Join(baseDir, schemaFile),
			filepath.Join(baseDir, viewFile),
			filepath.Join(baseDir, seedFile),
		),
		testcontainers.WithAdditionalWaitStrategyAndDeadline(
			containerStartTimeout,
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
			wait.ForListeningPort("5432/tcp"),
		),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start Postgres container")
	}
	return container, nil
}
