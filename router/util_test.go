package router_test

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/organization-service/goorg/router"
)

func setRouter(router router.IRouter, body string) {
	router.GET("/", write(body))
	router.POST("/", write(body))
	router.DELETE("/", write(body))
	router.PATCH("/", write(body))
	router.PUT("/", write(body))
}

func write(body string) func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		rw.Write([]byte(body))
	}
}
