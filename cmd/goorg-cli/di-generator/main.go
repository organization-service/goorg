package generator

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/organization-service/goorg/v2/cmd/goorg-cli/templates"
	"golang.org/x/tools/imports"
)

type config struct {
	SearchDir     string
	OutputDir     string
	DriverInclude bool
}

func dirwalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, dirwalk(filepath.Join(dir, file.Name()))...)
			continue
		}
		name := file.Name()
		pos := strings.LastIndex(name, ".")
		if pos > 0 && name[pos:] == ".go" {
			paths = append(paths, filepath.Join(dir, file.Name()))
		}
	}

	return paths
}

var keyMap = map[string]string{
	"infra.IDataBase": "di.DB",
}

func selectorExpr(se *ast.SelectorExpr) string {
	typ := ""
	name := se.Sel.Name
	switch xi := se.X.(type) {
	case *ast.Ident:
		typ = xi.Name
	}
	return typ + "." + name
}

func starExpr(n *ast.StarExpr) string {
	switch se := n.X.(type) {
	case *ast.SelectorExpr:
		return selectorExpr(se)
	}
	return ""
}

func returnType(fields *ast.FieldList) string {
	if fields != nil {
		for _, field := range fields.List {
			switch n := field.Type.(type) {
			case *ast.Ident:
				return n.Name
			case *ast.StarExpr:
				return starExpr(n)
			case *ast.SelectorExpr:
				return selectorExpr(n)
			default:
				return ""
			}
		}
	}
	return ""
}

// RegistryGenerator 登録
type RegistryGenerator struct {
}

type definition struct {
	FuncName string
}

// Definitions Definitions
type Definitions []*definition

// GetDefinition GetDefinition
func (r *RegistryGenerator) GetDefinition(path string) (ds Definitions, err error) {
	fs := token.NewFileSet()
	paths := dirwalk(path)
	for _, path := range paths {
		f, err := parser.ParseFile(fs, path, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		pack := f.Name.String() + "."
		for _, node := range f.Decls {
			switch n := node.(type) {
			case *ast.FuncDecl:
				comment := strings.TrimSpace(n.Doc.Text())
				switch {
				case strings.HasPrefix(comment, "@handler"), strings.HasPrefix(comment, "@di.h"):
					log.Println(path)
					r := returnType(n.Type.Results)
					consts := ""
					value := ""
					if r != "" {
						value = pack + r
						consts = r
						if _, ok := keyMap[value]; !ok {
							keyMap[value] = consts
						}
					}
					d := &definition{
						FuncName: pack + n.Name.String(),
					}
					ds = append(ds, d)
				case strings.HasPrefix(comment, "@service"), strings.HasPrefix(comment, "@usecase"), strings.HasPrefix(comment, "@di.s"), strings.HasPrefix(comment, "@di.u"):
					log.Println(path)
					r := returnType(n.Type.Results)
					consts := ""
					value := ""
					if r != "" {
						value = pack + r
						consts = r
						if _, ok := keyMap[value]; !ok {
							keyMap[value] = consts
						}
					}
					d := &definition{
						FuncName: pack + n.Name.String(),
					}
					ds = append(ds, d)
				case strings.HasPrefix(comment, "@repository"), strings.HasPrefix(comment, "@di.r"):
					log.Println(path)
					r := returnType(n.Type.Results)
					consts := ""
					value := ""
					if r != "" {
						value = r
						cs := strings.Split(r, ".")
						if len(cs) == 0 {
							return nil, errors.New("Dependency is not a domain driven design")
						}
						consts = cs[len(cs)-1]
						if _, ok := keyMap[value]; !ok {
							keyMap[value] = consts
						}
					}
					d := &definition{
						FuncName: pack + n.Name.String(),
					}
					ds = append(ds, d)
				}
			}
		}
	}
	return
}

// Execute テンプレートのパース処理
func (r *RegistryGenerator) Execute(templateStruct interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	funcParams := func(params []string) string {
		return "{" + strings.Join(params, ", ") + "}"
	}
	f := template.FuncMap{
		"params": funcParams,
	}
	stempl, err := templates.GetFile("di.tpl")
	if err != nil {
		return nil, err
	}
	tpl := template.Must(template.New("").Funcs(f).Parse(string(stempl)))
	if err := tpl.Execute(buf, templateStruct); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Import import処理
func (r *RegistryGenerator) Import(buf []byte) ([]byte, error) {
	buff, err := imports.Process("", buf, &imports.Options{
		FormatOnly: false,
		Comments:   true,
	})
	if err != nil {
		return nil, err
	}
	return buff, nil
}

// Build generate process
func (r *RegistryGenerator) Build(config *config) error {
	if _, err := os.Stat(config.SearchDir); os.IsNotExist(err) {
		return fmt.Errorf("dir: %s is not exist", config.SearchDir)
	}

	log.Println("Generate di-registry")

	type TemplateStruct struct {
		PackageName   string
		Definitions   Definitions
		Timestamp     string
		DriverInclude bool
	}
	absOutputDir, err := filepath.Abs(config.OutputDir)
	if err != nil {
		return err
	}
	packageName := filepath.Base(absOutputDir)

	templateStruct := TemplateStruct{
		PackageName:   packageName,
		Definitions:   Definitions{},
		Timestamp:     time.Now().Format("2006/01/02 15:04:05"),
		DriverInclude: config.DriverInclude,
	}
	templateStruct.Definitions, err = r.GetDefinition(config.SearchDir)
	if err != nil {
		return err
	}

	buff, err := r.Execute(templateStruct)
	if err != nil {
		return err
	}
	buff, err = r.Import(buff)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(config.OutputDir, os.ModePerm); err != nil {
		return err
	}
	registry, err := os.Create(filepath.Join(config.OutputDir, "registry.go"))
	if err != nil {
		return err
	}
	defer registry.Close()
	src, err := format.Source(buff)
	if err != nil {
		src = buff
	}
	_, err = registry.Write(src)
	return err
}

func Action(searchDirFlag, outputFlag string, driverInclude bool) error {
	g := RegistryGenerator{}
	conf := &config{
		SearchDir:     searchDirFlag,
		OutputDir:     outputFlag,
		DriverInclude: driverInclude,
	}
	return g.Build(conf)
}
