package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/organization-service/goorg/logger"
)

type (
	defaultRouter struct {
		*httprouter.Router
		group string
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
	r.Router.DELETE(joinURL(r, path), logger.Log(h))
}

// GET replaces httprouter.Router.GET.
func (r *defaultRouter) GET(path string, h httprouter.Handle) {
	r.Router.GET(joinURL(r, path), logger.Log(h))
}

// HEAD replaces httprouter.Router.HEAD.
func (r *defaultRouter) HEAD(path string, h httprouter.Handle) {
	r.Router.HEAD(joinURL(r, path), logger.Log(h))
}

// OPTIONS replaces httprouter.Router.OPTIONS.
func (r *defaultRouter) OPTIONS(path string, h httprouter.Handle) {
	r.Router.OPTIONS(joinURL(r, path), logger.Log(h))
}

// PATCH replaces httprouter.Router.PATCH.
func (r *defaultRouter) PATCH(path string, h httprouter.Handle) {
	r.Router.PATCH(joinURL(r, path), logger.Log(h))
}

// POST replaces httprouter.Router.POST.
func (r *defaultRouter) POST(path string, h httprouter.Handle) {
	r.Router.POST(joinURL(r, path), logger.Log(h))
}

// PUT replaces httprouter.Router.PUT.
func (r *defaultRouter) PUT(path string, h httprouter.Handle) {
	r.Router.PUT(joinURL(r, path), logger.Log(h))
}

// Handle replaces httprouter.Router.Handle.
func (r *defaultRouter) Handle(method, path string, h httprouter.Handle) {
	r.Router.Handle(method, joinURL(r, path), logger.Log(h))
}

// Handler replaces httprouter.Router.Handler.
func (r *defaultRouter) Handler(method, path string, handler http.Handler) {
	r.Router.Handler(method, joinURL(r, path), logger.LogHandler(handler))
}

// HandlerFunc replaces httprouter.Router.HandlerFunc.
func (r *defaultRouter) HandlerFunc(method, path string, handler http.HandlerFunc) {
	r.Router.HandlerFunc(method, joinURL(r, path), logger.LogHandlerFunc(handler))
}

// ServeHTTP replaces httprouter.Router.ServeHTTP.
func (r *defaultRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.Router.ServeHTTP(w, req)
}

func (r *defaultRouter) ServeFiles(path string, fileSystem http.FileSystem) {
	r.Router.ServeFiles(joinURL(r, path), fileSystem)
}

func (r *defaultRouter) GlobalOPTIONS(h http.HandlerFunc) {
	r.Router.GlobalOPTIONS = h
}

func (r *defaultRouter) Group(path string) IRouter {
	return &defaultRouter{
		Router: r.Router,
		group:  path,
	}
}

func (r *defaultRouter) getGroup() string {
	return r.group
}
