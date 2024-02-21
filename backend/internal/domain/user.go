package domain

import "github.com/oklog/ulid/v2"

// User ユーザー
type User struct {
	ID     ulid.ULID
	Name   string
	TeamID ulid.ULID
}
