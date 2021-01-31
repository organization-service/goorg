package generator_test

import (
	"log"
	"path/filepath"
	"testing"
	"time"

	generator "github.com/organization-service/goorg/cmd/di-generator"
)

func TestGen(t *testing.T) {
	path := "../../../examples"
	type TmplateStruct struct {
		PackageName string
		Definitions generator.Definitions
		Timestamp   string
	}
	absOutputDir, _ := filepath.Abs(path)
	packageName := filepath.Base(absOutputDir)
	r := generator.RegistryGenerator{}
	r.GetDefinition(path)
	loc, _ := time.LoadLocation("UTC")
	tim := time.Date(2020, 7, 24, 20, 0, 0, 0, loc)
	tmplateStruct := TmplateStruct{
		PackageName: packageName,
		Definitions: generator.Definitions{},
		Timestamp:   tim.Format("2006/01/02 15:04:05"),
	}
	tmplateStruct.Definitions, _ = r.GetDefinition(path)
	buff, _ := r.Execute(tmplateStruct)
	log.Println(string(buff))
}
