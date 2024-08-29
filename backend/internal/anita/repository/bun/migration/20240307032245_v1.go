package migration

import (
	"context"

	"github.com/ictsc/ictsc-outlands/backend/pkg/errors"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(
		func(ctx context.Context, db *bun.DB) error {
			_, err := db.NewCreateTable().Model((*v1Bastion)(nil)).Exec(ctx)
			if err != nil {
				return errors.Wrap(errors.ErrUnknown, err)
			}

			_, err = db.NewCreateTable().Model((*v1User)(nil)).Exec(ctx)
			if err != nil {
				return errors.Wrap(errors.ErrUnknown, err)
			}

			_, err = db.NewCreateTable().Model((*v1Team)(nil)).Exec(ctx)
			if err != nil {
				return errors.Wrap(errors.ErrUnknown, err)
			}

			return nil
		},
		func(ctx context.Context, db *bun.DB) error {
			_, err := db.NewDropTable().Model((*v1Team)(nil)).IfExists().Exec(ctx)
			if err != nil {
				return errors.Wrap(errors.ErrUnknown, err)
			}

			_, err = db.NewDropTable().Model((*v1User)(nil)).IfExists().Exec(ctx)
			if err != nil {
				return errors.Wrap(errors.ErrUnknown, err)
			}

			_, err = db.NewDropTable().Model((*v1Bastion)(nil)).IfExists().Exec(ctx)
			if err != nil {
				return errors.Wrap(errors.ErrUnknown, err)
			}

			return nil
		},
	)
}

// v1User ユーザーテーブル
type v1User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID     string `bun:"id,type:char(26),notnull,pk"`
	Name   string `bun:"name,type:varchar(20),notnull"`
	TeamID string `bun:"team_id,type:char(26),notnull"`

	Team *v1Team `bun:"rel:belongs-to,join:team_id=id"`
}

// v1Team チームテーブル
type v1Team struct {
	bun.BaseModel `bun:"table:teams,alias:t"`

	ID             string `bun:"id,type:char(26),notnull,pk"`
	Code           int    `bun:"code,type:tinyint,notnull,unique"`
	Name           string `bun:"name,type:varchar(20),notnull"`
	Organization   string `bun:"organization,type:varchar(50),notnull"`
	InvitationCode string `bun:"invitation_code,type:char(32),notnull,unique"`
	CodeRemaining  int    `bun:"codeRemaining,type:tinyint,notnull"`

	Bastion *v1Bastion `bun:"rel:has-one,join:id=team_id"`
	Members []*v1User  `bun:"rel:has-many,join:id=team_id"`
}

// v1Bastion 踏み台サーバー
type v1Bastion struct {
	bun.BaseModel `bun:"table:bastions,alias:b"`

	TeamID   string `bun:"team_id,type:char(26),notnull,pk"`
	User     string `bun:"user,type:varchar(20),notnull"`
	Password string `bun:"password,type:varchar(20),notnull"`
	Host     string `bun:"host,type:varchar(100),notnull,unique"`
	Port     int    `bun:"port,type:smallint unsigned,notnull"`
}
