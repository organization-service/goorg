package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type (
	IRouter interface {
		DELETE(path string, h httprouter.Handle)
		HEAD(path string, h httprouter.Handle)
		OPTIONS(path string, h httprouter.Handle)
		PATCH(path string, h httprouter.Handle)
		POST(path string, h httprouter.Handle)
		PUT(path string, h httprouter.Handle)
		GET(path string, h httprouter.Handle)
		Handle(method, path string, h httprouter.Handle)
		Handler(method, path string, handler http.Handler)
		HandlerFunc(method, path string, handler http.HandlerFunc)
		ServeHTTP(w http.ResponseWriter, req *http.Request)
		ServeFiles(path string, fileSystem http.FileSystem)
		GlobalOPTIONS(h http.HandlerFunc)
	}
)
