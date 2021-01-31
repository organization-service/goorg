package database

import (
	"context"
	"log"

	"github.com/organization-service/goorg/repository"
	"gorm.io/gorm"
)

type (
	Repository struct {
		driver IDriver
	}

	Connection struct {
		driver IDriver
	}

	Transaction struct {
		connectionReadWrite *gorm.DB
	}
)

func NewRepository(driver IDriver) *Repository {
	return &Repository{
		driver: driver,
	}
}

func (r *Repository) NewConnection() *Connection {
	return &Connection{
		driver: r.driver,
	}
}

func (con *Connection) Transaction(c context.Context, f func(tx repository.Transaction) error) error {
	var err error
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()
	tx := con.driver.ReadWriteConnection(c).Begin()
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
			panic(err)
		}
	}()
	err = f(&Transaction{
		connectionReadWrite: tx,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
