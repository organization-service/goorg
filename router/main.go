package router

import (
	"log"
	"net/http"
	"path"
	"runtime/debug"
	"strings"

	"github.com/organization-service/goorg/internal"
)

func New(fn ...func() interface{}) IRouter {
	var router IRouter
	switch strings.ToLower(internal.GetApmName()) {
	case internal.Elastic:
		router = newElastic(fn...)
	case internal.Newrelic:
		router = newNR()
	default:
		router = newDefault()
	}
	router.GlobalOPTIONS(globalOPTIONS())
	return router
}

func globalOPTIONS() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			header := w.Header()
			header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
			header.Set("Access-Control-Allow-Origin", "*")
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

func panicHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	log.Println(r.URL.Path, string(debug.Stack()))
	debug.PrintStack()
	w.WriteHeader(http.StatusInternalServerError)
}

func joinURL(router IRouter, url string) string {
	return path.Join(router.getGroup(), url)
}
