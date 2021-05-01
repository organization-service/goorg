package internal

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func CallHandler(h interface{}, w http.ResponseWriter, r *http.Request) {
	switch handler := h.(type) {
	case func(http.ResponseWriter, *http.Request, httprouter.Params):
		p := httprouter.ParamsFromContext(r.Context())
		handler(w, r, p)
	case httprouter.Handle:
		p := httprouter.ParamsFromContext(r.Context())
		handler(w, r, p)
	case http.HandlerFunc:
		handler.ServeHTTP(w, r)
	case http.Handler:
		handler.ServeHTTP(w, r)
	case func(http.ResponseWriter, *http.Request):
		handler(w, r)
	}
}
