package router

import (
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/newrelic/go-agent/v3/integrations/nrhttprouter"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/organization-service/goorg/logger"
)

type (
	nrRouter struct {
		*nrhttprouter.Router
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
	)
	return app
}

// DELETE replaces httprouter.Router.DELETE.
func (r *nrRouter) DELETE(path string, h httprouter.Handle) {
	r.Router.DELETE(path, logger.Log(h))
}

// GET replaces httprouter.Router.GET.
func (r *nrRouter) GET(path string, h httprouter.Handle) {
	r.Router.GET(path, logger.Log(h))
}

// HEAD replaces httprouter.Router.HEAD.
func (r *nrRouter) HEAD(path string, h httprouter.Handle) {
	r.Router.HEAD(path, logger.Log(h))
}

// OPTIONS replaces httprouter.Router.OPTIONS.
func (r *nrRouter) OPTIONS(path string, h httprouter.Handle) {
	r.Router.OPTIONS(path, logger.Log(h))
}

// PATCH replaces httprouter.Router.PATCH.
func (r *nrRouter) PATCH(path string, h httprouter.Handle) {
	r.Router.PATCH(path, logger.Log(h))
}

// POST replaces httprouter.Router.POST.
func (r *nrRouter) POST(path string, h httprouter.Handle) {
	r.Router.POST(path, logger.Log(h))
}

// PUT replaces httprouter.Router.PUT.
func (r *nrRouter) PUT(path string, h httprouter.Handle) {
	r.Router.PUT(path, logger.Log(h))
}

// Handle replaces httprouter.Router.Handle.
func (r *nrRouter) Handle(method, path string, h httprouter.Handle) {
	r.Router.Handle(method, path, logger.Log(h))
}

// Handler replaces httprouter.Router.Handler.
func (r *nrRouter) Handler(method, path string, handler http.Handler) {
	r.Router.Handler(method, path, logger.LogHandler(handler))
}

// HandlerFunc replaces httprouter.Router.HandlerFunc.
func (r *nrRouter) HandlerFunc(method, path string, handler http.HandlerFunc) {
	r.Router.HandlerFunc(method, path, logger.LogHandlerFunc(handler))
}

// ServeHTTP replaces httprouter.Router.ServeHTTP.
func (r *nrRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.Router.ServeHTTP(w, req)
}

func (r *nrRouter) ServeFiles(path string, fileSystem http.FileSystem) {
	r.Router.ServeFiles(path, fileSystem)
}

func (r *nrRouter) GlobalOPTIONS(h http.HandlerFunc) {
	r.Router.GlobalOPTIONS = h
}
