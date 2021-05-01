package router

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"runtime/debug"

	"github.com/julienschmidt/httprouter"
	"github.com/organization-service/goorg/v2/internal"
)

func New(fn ...func() interface{}) IRouter {
	var router IRouter
	switch internal.GetApmName() {
	case internal.Elastic:
		router = newElastic(fn...)
	case internal.Newrelic:
		router = newNR()
	default:
		router = newDefault()
	}
	router.GlobalOPTIONS(globalOPTIONS())
	healthCheck(router)
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

func healthCheck(router IRouter) {
	router.GET("/health", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(fmt.Sprintf(`{"status": "%s"}`, http.StatusText(http.StatusOK))))
	})
}
