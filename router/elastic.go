package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/organization-service/goorg/logger"
	"go.elastic.co/apm/module/apmhttprouter"
)

type (
	elasticRouter struct {
		*apmhttprouter.Router
	}
)

func newElastic(fn ...func() interface{}) IRouter {
	o := []apmhttprouter.Option{}
	if len(fn) > 0 {
		o = fn[0]().([]apmhttprouter.Option)
	}
	router := &elasticRouter{
		Router: apmhttprouter.New(o...),
	}
	router.PanicHandler = panicHandler
	return router
}

// DELETE replaces httprouter.Router.DELETE.
func (r *elasticRouter) DELETE(path string, h httprouter.Handle) {
	r.Router.DELETE(path, logger.Log(h))
}

// GET replaces httprouter.Router.GET.
func (r *elasticRouter) GET(path string, h httprouter.Handle) {
	r.Router.GET(path, logger.Log(h))
}

// HEAD replaces httprouter.Router.HEAD.
func (r *elasticRouter) HEAD(path string, h httprouter.Handle) {
	r.Router.HEAD(path, logger.Log(h))
}

// OPTIONS replaces httprouter.Router.OPTIONS.
func (r *elasticRouter) OPTIONS(path string, h httprouter.Handle) {
	r.Router.OPTIONS(path, logger.Log(h))
}

// PATCH replaces httprouter.Router.PATCH.
func (r *elasticRouter) PATCH(path string, h httprouter.Handle) {
	r.Router.PATCH(path, logger.Log(h))
}

// POST replaces httprouter.Router.POST.
func (r *elasticRouter) POST(path string, h httprouter.Handle) {
	r.Router.POST(path, logger.Log(h))
}

// PUT replaces httprouter.Router.PUT.
func (r *elasticRouter) PUT(path string, h httprouter.Handle) {
	r.Router.PUT(path, logger.Log(h))
}

// Handle replaces httprouter.Router.Handle.
func (r *elasticRouter) Handle(method, path string, h httprouter.Handle) {
	r.Router.Handle(method, path, logger.Log(h))
}

// Handler replaces httprouter.Router.Handler.
func (r *elasticRouter) Handler(method, path string, handler http.Handler) {
	r.Router.Handler(method, path, logger.LogHandler(handler))
}

// HandlerFunc replaces httprouter.Router.HandlerFunc.
func (r *elasticRouter) HandlerFunc(method, path string, handler http.HandlerFunc) {
	r.Router.HandlerFunc(method, path, logger.LogHandlerFunc(handler))
}

// ServeHTTP replaces httprouter.Router.ServeHTTP.
func (r *elasticRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.Router.ServeHTTP(w, req)
}

func (r *elasticRouter) ServeFiles(path string, fileSystem http.FileSystem) {
	r.Router.ServeFiles(path, fileSystem)
}

func (r *elasticRouter) GlobalOPTIONS(h http.HandlerFunc) {
	r.Router.GlobalOPTIONS = h
}
