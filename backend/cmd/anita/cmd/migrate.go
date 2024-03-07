package cmd

import (
	"log"

	"github.com/ictsc/ictsc-outlands/backend/internal/anita/repository/bun/migration"
	"github.com/ictsc/ictsc-outlands/backend/pkg/db/rdb/bun"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun/migrate"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Execute bun migration",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := bun.NewDB(provideRDBConfig(&config))
		if err != nil {
			log.Panic(err)
		}
		bunDB := db.GetBunDB()

		migrator := migrate.NewMigrator(bunDB, migration.Migrations)
		ctx := cmd.Context()

		if err = migrator.Init(ctx); err != nil {
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

func init() {
	rootCmd.AddCommand(migrateCmd)
}
