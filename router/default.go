package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/organization-service/goorg/logger"
)

type (
	defaultRouter struct {
		*httprouter.Router
	}
)

func newDefault() IRouter {
	router := &defaultRouter{
		Router: httprouter.New(),
	}
	router.PanicHandler = panicHandler
	return router
}

// DELETE replaces httprouter.Router.DELETE.
func (r *defaultRouter) DELETE(path string, h httprouter.Handle) {
	r.Router.DELETE(path, logger.Log(h))
}

// GET replaces httprouter.Router.GET.
func (r *defaultRouter) GET(path string, h httprouter.Handle) {
	r.Router.GET(path, logger.Log(h))
}

// HEAD replaces httprouter.Router.HEAD.
func (r *defaultRouter) HEAD(path string, h httprouter.Handle) {
	r.Router.HEAD(path, logger.Log(h))
}

// OPTIONS replaces httprouter.Router.OPTIONS.
func (r *defaultRouter) OPTIONS(path string, h httprouter.Handle) {
	r.Router.OPTIONS(path, logger.Log(h))
}

// PATCH replaces httprouter.Router.PATCH.
func (r *defaultRouter) PATCH(path string, h httprouter.Handle) {
	r.Router.PATCH(path, logger.Log(h))
}

// POST replaces httprouter.Router.POST.
func (r *defaultRouter) POST(path string, h httprouter.Handle) {
	r.Router.POST(path, logger.Log(h))
}

// PUT replaces httprouter.Router.PUT.
func (r *defaultRouter) PUT(path string, h httprouter.Handle) {
	r.Router.PUT(path, logger.Log(h))
}

// Handle replaces httprouter.Router.Handle.
func (r *defaultRouter) Handle(method, path string, h httprouter.Handle) {
	r.Router.Handle(method, path, logger.Log(h))
}

// Handler replaces httprouter.Router.Handler.
func (r *defaultRouter) Handler(method, path string, handler http.Handler) {
	r.Router.Handler(method, path, logger.LogHandler(handler))
}

// HandlerFunc replaces httprouter.Router.HandlerFunc.
func (r *defaultRouter) HandlerFunc(method, path string, handler http.HandlerFunc) {
	r.Router.HandlerFunc(method, path, logger.LogHandlerFunc(handler))
}

// ServeHTTP replaces httprouter.Router.ServeHTTP.
func (r *defaultRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.Router.ServeHTTP(w, req)
}

func (r *defaultRouter) ServeFiles(path string, fileSystem http.FileSystem) {
	r.Router.ServeFiles(path, fileSystem)
}

func (r *defaultRouter) GlobalOPTIONS(h http.HandlerFunc) {
	r.Router.GlobalOPTIONS = h
}
