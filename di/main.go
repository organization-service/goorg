package di

import (
	"go.uber.org/dig"
)

type (
	Container struct {
		container *dig.Container
		err       error
	}
)

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
