// template file generate command: statik -src=templates -f

package generator

import (
	"bytes"
	"go/format"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/rakyll/statik/fs"
	"golang.org/x/tools/imports"
)

type (
	templateModel struct {
		Timestamp string
	}
)

func newTemplateModel() *templateModel {
	return &templateModel{
		Timestamp: time.Now().Format("2006/01/02 15:04:05"),
	}
}

func importModule(buf []byte) ([]byte, error) {
	buff, err := imports.Process("", buf, &imports.Options{
		FormatOnly: false,
		Comments:   true,
	})
	if err != nil {
		return nil, err
	}
	return buff, nil
}

func execute(tmplateStruct interface{}, fileName string) ([]byte, error) {
	fileSystem, err := fs.New()
	if err != nil {
		return nil, err
	}
	f, err := fileSystem.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tmpBuf, _ := ioutil.ReadAll(f)
	buf := new(bytes.Buffer)
	tmpl := template.Must(template.New("").Parse(string(tmpBuf)))
	if err := tmpl.Execute(buf, tmplateStruct); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func tmpToBuffer(tmpFileName string, tmpModel *templateModel) ([]byte, error) {
	var buf []byte
	var err error
	if buf, err = execute(tmpModel, tmpFileName); err != nil {
		return nil, err
	}
	if buf, err = importModule(buf); err != nil {
		return nil, err
	}

	return buf, nil
}

func fileCreate(buf []byte, outputDir, outputFileName string) error {
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return err
	}
	registry, err := os.Create(filepath.Join(outputDir, outputFileName))
	if err != nil {
		return err
	}
	defer registry.Close()
	src, err := format.Source(buf)
	if err != nil {
		src = buf
	}
	_, err = registry.Write(src)
	return err
}

func build() error {
	log.Println("Generate swagger file")
	tmpModel := newTemplateModel()
	buf, err := tmpToBuffer("/swagger.tpl", tmpModel)
	if err != nil {
		return err
	}
	if err := fileCreate(buf, "application/server", "swagger.go"); err != nil {
		return err
	}
	return nil
}

func Action() error {
	return build()
}
