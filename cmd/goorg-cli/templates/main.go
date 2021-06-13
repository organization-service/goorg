package templates

import (
	"embed"
	"io"
	"io/fs"
)

//go:embed templates/*
var assets embed.FS

var (
	template fs.FS
)

func init() {
	var err error
	template, err = fs.Sub(assets, "templates")
	if err != nil {
		panic(err)
	}
}

func GetFile(filename string) ([]byte, error) {
	f, err := template.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
	// // filename = filepath.Join("templates", filename)
	// return fs.ReadFile(assets, filename)
}
