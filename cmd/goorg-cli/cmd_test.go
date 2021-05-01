package main_test

import (
	"testing"

	swaggerGene "github.com/organization-service/goorg/v2/cmd/goorg-cli/swagger-generator"
)

func TestSwagger(t *testing.T) {
	swaggerGene.Action()
}
