package infra

import (
	"context"
	"log"

	"gorm.io/gorm"
)

type (
	dbRepository struct {
		driver IDriver
	}

	dbConnection struct {
		driver IDriver
	}

	dbTransaction struct {
		master *gorm.DB
		slave  *gorm.DB
	}
)

func NewRepository(driver IDriver) repository.Repository {
	return &dbRepository{
		driver: driver,
	}
}

func (r *dbRepository) NewConnection() (repository.Connection, error) {
	return &dbConnection{
		driver: r.driver,
	}, nil
}

func (r *dbRepository) MustConnection() repository.Connection {
	db, err := r.NewConnection()
	if err != nil {
		panic(err)
	}
	return db
}

func (con *dbConnection) Transaction(c context.Context, f func(tx repository.Transaction) error) error {
	var err error
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()
	tx := con.driver.Master(c).Begin()
	err = f(&dbTransaction{
		master: tx,
		slave:  con.driver.Slave(c),
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
