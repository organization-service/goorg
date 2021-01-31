package infra

import (
	"github.com/organization-service/goorg/database"
	"github.com/organization-service/goorg/repository"
)

type (
	dbRepository struct {
		repo *database.Repository
	}

	dbConnection struct {
		*database.Connection
	}

	dbTransaction struct {
		*database.Transaction
	}
)

// @repository
func NewRepository(driver database.IDriver) repository.Repository {
	return &dbRepository{
		repo: database.NewRepository(driver),
	}
}

func (r *dbRepository) NewConnection() (repository.Connection, error) {
	return &dbConnection{
		Connection: r.repo.NewConnection(),
	}, nil
}

func (r *dbRepository) MustConnection() repository.Connection {
	db, err := r.NewConnection()
	if err != nil {
		panic(err)
	}
	return db
}
