package repository

import (
	"github.com/organization-service/goorg/repository"
)

type (
	Repository interface {
		NewConnection() (Connection, error)
		MustConnection() Connection
	}
	Connection interface {
		repository.Connection
	}
	Transaction interface {
	}
)
