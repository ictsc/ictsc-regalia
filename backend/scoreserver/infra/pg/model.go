package pg

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Team struct {
	ID           pgtype.UUID `db:"id"`
	Code         int64       `db:"code"`
	Name         string      `db:"name"`
	Organization string      `db:"organization"`
	CreatedAt    time.Time   `db:"created_at"`
	UpdateAt     time.Time   `db:"updated_at"`
}
