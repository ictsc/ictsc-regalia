// Package bun Bunユーティリティー
package bun

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/ictsc/ictsc-outlands/backend/pkg/db/rdb"
	"github.com/ictsc/ictsc-outlands/backend/pkg/errors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/extra/bundebug"
)

// DB Bun RDBクライアント
type DB struct {
	*bun.DB
}

var _ rdb.Tx = (*DB)(nil)

// RDB接続設定
type Config struct {
	Dev bool

	Hostname string
	Port     int
	Username string
	Password string
	Database string
}

// NewDB Bun RDBクライアント生成
func NewDB(conf *Config) (*DB, error) {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, errors.Wrap(err)
	}

	mysqlConf := mysql.NewConfig()
	mysqlConf.User = conf.Username
	mysqlConf.Passwd = conf.Password
	mysqlConf.Addr = conf.Hostname + ":" + strconv.Itoa(conf.Port)
	mysqlConf.DBName = conf.Database
	mysqlConf.ParseTime = true
	mysqlConf.Loc = jst

	sqlDB, err := sql.Open("mysql", mysqlConf.FormatDSN())
	if err != nil {
		return nil, errors.Wrap(err)
	}

	bunDB := bun.NewDB(sqlDB, mysqldialect.New())
	if conf.Dev {
		bunDB.AddQueryHook(bundebug.NewQueryHook())
	}

	return &DB{DB: bunDB}, nil
}
