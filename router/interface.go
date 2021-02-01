package router

import (
	"net/http"
)

type (
	IRouter interface {
		DELETE(path string, h interface{})
		HEAD(path string, h interface{})
		OPTIONS(path string, h interface{})
		PATCH(path string, h interface{})
		POST(path string, h interface{})
		PUT(path string, h interface{})
		GET(path string, h interface{})
		Handle(method, path string, h interface{})
		Handler(method, path string, h http.Handler)
		HandlerFunc(method, path string, h http.HandlerFunc)
		ServeHTTP(w http.ResponseWriter, req *http.Request)
		ServeFiles(path string, fileSystem http.FileSystem)
		GlobalOPTIONS(h http.HandlerFunc)
		Group(path string) IRouter
		getGroup() string
	}
)
