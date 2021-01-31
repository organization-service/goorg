package di

import (
	"github.com/organization-service/goorg/database"
	"github.com/organization-service/goorg/repository"
	"go.uber.org/dig"
)

type (
	Container struct {
		container *dig.Container
		err       error
	}
	repoFunc func(driver database.IDriver) repository.Repository
)

func (di *Container) Driver(driver ...func() database.IDriver) *Container {
	value := database.New
	if len(driver) > 0 {
		value = driver[0]
	}
	di.err = di.container.Provide(value)
	return di
}

func (di *Container) Repository(repo ...repoFunc) *Container {
	value := database.NewRepository
	if len(repo) > 0 {
		value = repo[0]
	}
	di.err = di.container.Provide(value)
	return di
}

func (di *Container) Provide(constructor interface{}, opts ...dig.ProvideOption) *Container {
	di.err = di.container.Provide(constructor, opts...)
	return di
}

func (di *Container) Invoke(function interface{}, opts ...dig.InvokeOption) *Container {
	di.err = di.container.Invoke(function, opts...)
	return di
}

func (di *Container) Error() error {
	return di.err
}

func New() *Container {
	container := dig.New()
	di := &Container{
		container: container,
	}
	return di
}
