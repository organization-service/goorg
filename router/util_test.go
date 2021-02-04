package router_test

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/organization-service/goorg/router"
)

var (
	sleepTime time.Duration = 0
)

func setRouter(router router.IRouter, body string) {
	router.GET("/", write(body))
	router.POST("/", write(body))
	router.DELETE("/", write(body))
	router.PATCH("/", write(body))
	router.PUT("/", write(body))
	{
		router := router.Group("/api")
		{
			router := router.Group("/test")
			{
				router.GET("", write(body))
			}
		}
	}
	router.HandlerFunc(http.MethodGet, "/handler-func", writeHandler(body))
	router.Handler(http.MethodGet, "/handler", http.HandlerFunc(writeHandlerFunc(body)))
}

func setRouterHandlerFunc(router router.IRouter, body string) {
	router.GET("/", writeHandlerFunc(body))
	router.POST("/", writeHandlerFunc(body))
	router.DELETE("/", writeHandlerFunc(body))
	router.PATCH("/", writeHandlerFunc(body))
	router.PUT("/", writeHandlerFunc(body))
	{
		router := router.Group("/api")
		{
			router := router.Group("/test")
			{
				router.GET("", writeHandlerFunc(body))
			}
		}
	}
	router.HandlerFunc(http.MethodGet, "/handler-func", writeHandler(body))
	router.Handler(http.MethodGet, "/handler", http.HandlerFunc(writeHandlerFunc(body)))
}

func setRouterHandler(router router.IRouter, body string) {
	router.GET("/", writeHandler(body))
	router.POST("/", writeHandler(body))
	router.DELETE("/", writeHandler(body))
	router.PATCH("/", writeHandler(body))
	router.PUT("/", writeHandler(body))
	{
		router := router.Group("/api")
		{
			router := router.Group("/test")
			{
				router.GET("", writeHandler(body))
			}
		}
	}
	router.HandlerFunc(http.MethodGet, "/handler-func", writeHandler(body))
	router.Handler(http.MethodGet, "/handler", http.HandlerFunc(writeHandlerFunc(body)))
}

func writeHandler(body string) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		time.Sleep(sleepTime * time.Second)
		rw.Write([]byte(body))
	}
}

func writeHandlerFunc(body string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		time.Sleep(sleepTime * time.Second)
		rw.Write([]byte(body))
	}
}

func write(body string) func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		time.Sleep(sleepTime * time.Second)
		rw.Write([]byte(body))
	}
}
