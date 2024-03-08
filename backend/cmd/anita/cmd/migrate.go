package cmd

import (
	"context"
	"log"

	"github.com/ictsc/ictsc-outlands/backend/internal/anita/repository/bun/migration"
	"github.com/ictsc/ictsc-outlands/backend/pkg/db/rdb/bun"
	"github.com/ictsc/ictsc-outlands/backend/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun/migrate"
)

func setupMigrator(ctx context.Context, config *Config) (*migrate.Migrator, error) {
	db, err := bun.NewDB(provideRDBConfig(config))
	if err != nil {
		return nil, err
	}

	bunDB := db.GetBunDB()
	migrator := migrate.NewMigrator(bunDB, migration.Migrations)

	if err = migrator.Init(ctx); err != nil {
		return nil, errors.Wrap(errors.ErrUnknown, err)
	}

	return migrator, nil
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Execute bun migration",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		migrator, err := setupMigrator(ctx, &config)
		if err != nil {
			log.Panic(err)
		}

		if err = migrator.Lock(ctx); err != nil {
			log.Panic(err)
		}
		defer func() {
			if err = migrator.Unlock(ctx); err != nil {
				log.Panic(err)
			}
		}()

		group, err := migrator.Migrate(ctx)
		if err != nil {
			log.Panic(err)
		}

		if group.IsZero() {
			log.Println("No changes")
		} else {
			log.Printf("Migrated to %s\n", group)
		}
	},
}

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Execute bun migration rollback",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		migrator, err := setupMigrator(ctx, &config)
		if err != nil {
			log.Panic(err)
		}

		if err = migrator.Lock(ctx); err != nil {
			log.Panic(err)
		}
		defer func() {
			if err = migrator.Unlock(ctx); err != nil {
				log.Panic(err)
			}
		}()

		group, err := migrator.Rollback(ctx)
		if err != nil {
			log.Panic(err)
		}

		if group.IsZero() {
			log.Println("No changes")
		} else {
			log.Printf("Rolled back to %s\n", group)
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(rollbackCmd)
}
