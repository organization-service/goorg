package database

import (
	"context"
	"database/sql"
	"errors"

	"gorm.io/gorm"
)

type (
	IDriver interface {
		Master(c context.Context) *gorm.DB
		MasterSql(c context.Context) *sql.DB
		Slave(c context.Context) *gorm.DB
		SlaveSql(c context.Context) *sql.DB
		LogMode(set bool)
	}
)

var (
	configs        map[string]*environment
	errNotFoundKey         = errors.New("Not found key")
	Master                 = "master"
	Slave                  = "slave"
	db             IDriver = nil
)
