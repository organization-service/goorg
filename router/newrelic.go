package router

import (
	"net/http"
	"os"
	"path"

	"github.com/julienschmidt/httprouter"
	"github.com/newrelic/go-agent/v3/integrations/nrhttprouter"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type (
	nrRouter struct {
		*nrhttprouter.Router
		group string
	}
)

func newNR() IRouter {
	app := newrelicApplication()
	router := &nrRouter{
		Router: nrhttprouter.New(app),
	}
	router.PanicHandler = panicHandler
	return router
}

// newrelicApplication newrelic applicationの設定を行うメソッド
func newrelicApplication() *newrelic.Application {
	NEW_RELIC_APP_NAME := os.Getenv("NEW_RELIC_APP_NAME")
	NEW_RELIC_LICENSE_KEY := os.Getenv("NEW_RELIC_LICENSE_KEY")
	app, _ := newrelic.NewApplication(
		newrelic.ConfigAppName(NEW_RELIC_APP_NAME),
		newrelic.ConfigLicense(NEW_RELIC_LICENSE_KEY),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	return app
}

// DELETE replaces httprouter.Router.DELETE.
func (r *nrRouter) DELETE(path string, h interface{}) {
	r.Router.DELETE(joinURL(r, path), logHandler(h).(func(w http.ResponseWriter, r *http.Request, p httprouter.Params)))
}

// GET replaces httprouter.Router.GET.
func (r *nrRouter) GET(path string, h interface{}) {
	r.Router.GET(joinURL(r, path), logHandler(h).(func(w http.ResponseWriter, r *http.Request, p httprouter.Params)))
}

// HEAD replaces httprouter.Router.HEAD.
func (r *nrRouter) HEAD(path string, h interface{}) {
	r.Router.HEAD(joinURL(r, path), logHandler(h).(func(w http.ResponseWriter, r *http.Request, p httprouter.Params)))
}

// OPTIONS replaces httprouter.Router.OPTIONS.
func (r *nrRouter) OPTIONS(path string, h interface{}) {
	r.Router.OPTIONS(joinURL(r, path), logHandler(h).(func(w http.ResponseWriter, r *http.Request, p httprouter.Params)))
}

// PATCH replaces httprouter.Router.PATCH.
func (r *nrRouter) PATCH(path string, h interface{}) {
	r.Router.PATCH(joinURL(r, path), logHandler(h).(func(w http.ResponseWriter, r *http.Request, p httprouter.Params)))
}

// POST replaces httprouter.Router.POST.
func (r *nrRouter) POST(path string, h interface{}) {
	r.Router.POST(joinURL(r, path), logHandler(h).(func(w http.ResponseWriter, r *http.Request, p httprouter.Params)))
}

// PUT replaces httprouter.Router.PUT.
func (r *nrRouter) PUT(path string, h interface{}) {
	r.Router.PUT(joinURL(r, path), logHandler(h).(func(w http.ResponseWriter, r *http.Request, p httprouter.Params)))
}

// Handle replaces httprouter.Router.Handle.
func (r *nrRouter) Handle(method, path string, h interface{}) {
	r.Router.Handle(method, joinURL(r, path), logHandler(h).(func(w http.ResponseWriter, r *http.Request, p httprouter.Params)))
}

// Handler replaces httprouter.Router.Handler.
func (r *nrRouter) Handler(method, path string, h http.Handler) {
	r.Router.Handler(method, joinURL(r, path), logHandler(h, handler).(http.Handler))
}

// HandlerFunc replaces httprouter.Router.HandlerFunc.
func (r *nrRouter) HandlerFunc(method, path string, h http.HandlerFunc) {
	r.Router.HandlerFunc(method, joinURL(r, path), logHandler(h, handlerFunc).(http.HandlerFunc))
}

// ServeHTTP replaces httprouter.Router.ServeHTTP.
func (r *nrRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.Router.ServeHTTP(w, req)
}

func (r *nrRouter) ServeFiles(path string, fileSystem http.FileSystem) {
	r.Router.ServeFiles(joinURL(r, path), fileSystem)
}

func (r *nrRouter) GlobalOPTIONS(h http.HandlerFunc) {
	r.Router.GlobalOPTIONS = h
}

func (r *nrRouter) Group(url string) IRouter {
	return &nrRouter{
		Router: r.Router,
		group:  path.Join(r.group, url),
	}
}

func (r *nrRouter) getGroup() string {
	return r.group
}
