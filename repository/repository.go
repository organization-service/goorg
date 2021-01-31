package repository

import "context"

type (
	Repository interface {
		NewConnection() (Connection, error)
		MustConnection() Connection
	}
	Connection interface {
		Transaction(c context.Context, f func(tx interface{}) error) error
	}
	Transaction interface {
	}
)
