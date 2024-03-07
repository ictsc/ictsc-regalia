package bun

import (
	"github.com/uptrace/bun"
)

// User ユーザーテーブル
type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID     string `bun:"id,type:char(26),notnull,pk"`
	Name   string `bun:"name,type:varchar(20),notnull"`
	TeamID string `bun:"team_id,type:char(26),notnull"`

	Team *Team `bun:"rel:belongs-to,join:team_id=id"`
}

// Team チームテーブル
type Team struct {
	bun.BaseModel `bun:"table:teams,alias:t"`

	ID             string `bun:"id,type:char(26),notnull,pk"`
	Code           int    `bun:"code,type:tinyint,notnull,unique"`
	Name           string `bun:"name,type:varchar(20),notnull"`
	Organization   string `bun:"organization,type:varchar(50),notnull"`
	InvitationCode string `bun:"invitation_code,type:char(32),notnull,unique"`
	CodeRemaining  int    `bun:"codeRemaining,type:tinyint,notnull"`

	Bastion *Bastion `bun:"rel:has-one,join:id=team_id"`
	Members []*User  `bun:"rel:has-many,join:id=team_id"`
}

// Bastion 踏み台サーバー
type Bastion struct {
	bun.BaseModel `bun:"table:bastions,alias:b"`

	TeamID   string `bun:"team_id,type:char(26),notnull,pk"`
	User     string `bun:"user,type:varchar(20),notnull"`
	Password string `bun:"password,type:varchar(20),notnull"`
	Host     string `bun:"host,type:varchar(100),notnull,unique"`
	Port     int    `bun:"port,type:smallint unsigned,notnull"`
}
