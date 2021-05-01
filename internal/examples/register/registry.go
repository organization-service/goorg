// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by goorg/di/generator at 2021/01/31 22:55:34

package register

import (
	"github.com/organization-service/goorg/v2/di"
	"github.com/organization-service/goorg/v2/internal/examples/application/server"
	"github.com/organization-service/goorg/v2/internal/examples/infra"
	"github.com/organization-service/goorg/v2/internal/examples/usecase"
)

func New() *di.Container {
	container := di.New()
	container.Provide(server.NewHandler)
	container.Provide(infra.NewRepository)
	container.Provide(usecase.NewUseCase)

	return container
}
