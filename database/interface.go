package database

import (
	"context"
	"database/sql"
	"errors"

	"gorm.io/gorm"
)

type (
	IDriver interface {
		ReadWriteConnection(c context.Context) *gorm.DB
		ReadWriteConnectionObject(c context.Context) *sql.DB
		ReadOnlyConnection(c context.Context) *gorm.DB
		ReadOnlyConnectionObject(c context.Context) *sql.DB
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
