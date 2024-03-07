package migration

import (
	"github.com/uptrace/bun/migrate"
)

// Migrations マイグレーション一覧
var Migrations = migrate.NewMigrations() // nolint:gochecknoglobals
