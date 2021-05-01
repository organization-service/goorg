package router

import (
	"net/http"
	"path"

	"github.com/julienschmidt/httprouter"
	"go.elastic.co/apm/module/apmhttprouter"
)

type (
	elasticRouter struct {
		*apmhttprouter.Router
		group string
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
func (r *elasticRouter) DELETE(path string, h interface{}) {
	r.Router.DELETE(joinURL(r, path), logHandler(h).(func(w http.ResponseWriter, r *http.Request, p httprouter.Params)))
}

// GET replaces httprouter.Router.GET.
func (r *elasticRouter) GET(path string, h interface{}) {
	r.Router.GET(joinURL(r, path), logHandler(h).(func(w http.ResponseWriter, r *http.Request, p httprouter.Params)))
}

// HEAD replaces httprouter.Router.HEAD.
func (r *elasticRouter) HEAD(path string, h interface{}) {
	r.Router.HEAD(joinURL(r, path), logHandler(h).(func(w http.ResponseWriter, r *http.Request, p httprouter.Params)))
}

// OPTIONS replaces httprouter.Router.OPTIONS.
func (r *elasticRouter) OPTIONS(path string, h interface{}) {
	r.Router.OPTIONS(joinURL(r, path), logHandler(h).(func(w http.ResponseWriter, r *http.Request, p httprouter.Params)))
}

// PATCH replaces httprouter.Router.PATCH.
func (r *elasticRouter) PATCH(path string, h interface{}) {
	r.Router.PATCH(joinURL(r, path), logHandler(h).(func(w http.ResponseWriter, r *http.Request, p httprouter.Params)))
}

// POST replaces httprouter.Router.POST.
func (r *elasticRouter) POST(path string, h interface{}) {
	r.Router.POST(joinURL(r, path), logHandler(h).(func(w http.ResponseWriter, r *http.Request, p httprouter.Params)))
}

// PUT replaces httprouter.Router.PUT.
func (r *elasticRouter) PUT(path string, h interface{}) {
	r.Router.PUT(joinURL(r, path), logHandler(h).(func(w http.ResponseWriter, r *http.Request, p httprouter.Params)))
}

// Handle replaces httprouter.Router.Handle.
func (r *elasticRouter) Handle(method, path string, h interface{}) {
	r.Router.Handle(method, joinURL(r, path), logHandler(h).(func(w http.ResponseWriter, r *http.Request, p httprouter.Params)))
}

// Handler replaces httprouter.Router.Handler.
func (r *elasticRouter) Handler(method, path string, h http.Handler) {
	r.Router.Handler(method, joinURL(r, path), logHandler(h, handler).(http.Handler))
}

// HandlerFunc replaces httprouter.Router.HandlerFunc.
func (r *elasticRouter) HandlerFunc(method, path string, h http.HandlerFunc) {
	r.Router.HandlerFunc(method, joinURL(r, path), logHandler(h, handlerFunc).(http.HandlerFunc))
}

// ServeHTTP replaces httprouter.Router.ServeHTTP.
func (r *elasticRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.Router.ServeHTTP(w, req)
}

func (r *elasticRouter) ServeFiles(path string, fileSystem http.FileSystem) {
	r.Router.ServeFiles(joinURL(r, path), fileSystem)
}

func (r *elasticRouter) GlobalOPTIONS(h http.HandlerFunc) {
	r.Router.GlobalOPTIONS = h
}

func (r *elasticRouter) Group(url string) IRouter {
	return &elasticRouter{
		Router: r.Router,
		group:  path.Join(r.group, url),
	}
}

func (r *elasticRouter) getGroup() string {
	return r.group
}
