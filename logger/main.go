package logger

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type (
	Response struct {
		StatusCode int
		Header     http.Header
	}
	responseWriter struct {
		http.ResponseWriter
		resp Response
	}
)

func (w *responseWriter) Write(buf []byte) (int, error) {
	n, e := w.ResponseWriter.Write(buf)
	if w.resp.StatusCode == 0 {
		w.resp.StatusCode = http.StatusOK
	}
	return n, e
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.resp.StatusCode = statusCode
}

func (w *responseWriter) CloseNotify() <-chan bool {
	if closeNotifier, ok := w.ResponseWriter.(http.CloseNotifier); ok {
		return closeNotifier.CloseNotify()
	}
	return nil
}

func (w *responseWriter) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func WrapResponseWriter(w http.ResponseWriter) (http.ResponseWriter, *Response) {
	rw := responseWriter{
		ResponseWriter: w,
		resp: Response{
			Header: w.Header(),
		},
	}

	return &rw, &rw.resp
}

func Log(h interface{}) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		rAddr := r.Header.Get("X-Real-IP")
		if rAddr == "" {
			rAddr = r.RemoteAddr
		}
		method := r.Method
		path := r.URL.Path
		log.Printf("Remote:[%-20.20s]:[%-6.6s]:[%-50.50s]", rAddr, method, path)
		w, resp := WrapResponseWriter(w)
		if r.Method == http.MethodOptions {
			return
		}
		switch handler := h.(type) {
		case func(http.ResponseWriter, *http.Request, httprouter.Params):
			handler(w, r, p)
		case httprouter.Handle:
			handler(w, r, p)
		case http.HandlerFunc:
			*r = *r.WithContext(context.WithValue(r.Context(), httprouter.ParamsKey, p))
			handler.ServeHTTP(w, r)
		case func(http.ResponseWriter, *http.Request):
			*r = *r.WithContext(context.WithValue(r.Context(), httprouter.ParamsKey, p))
			handler(w, r)
		default:
			panic(errors.New("Not type handler"))
		}
		log.Println(fmt.Sprintf("Remote:[%-20.20s]:[%-6.6s]:[%-50.50s]:Status:[%d]", rAddr, method, path, resp.StatusCode))
	}
}

func LogHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rAddr := r.RemoteAddr
		method := r.Method
		path := r.URL.Path
		log.Printf("Remote: %s [%s] %s", rAddr, method, path)
		w, resp := WrapResponseWriter(w)
		if r.Method == http.MethodOptions {
			return
		}
		h.ServeHTTP(w, r)
		log.Println(fmt.Sprintf("Status: %v", resp.StatusCode))
	})
}

func LogHandlerFunc(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rAddr := r.RemoteAddr
		method := r.Method
		path := r.URL.Path
		log.Printf("Remote: %s [%s] %s", rAddr, method, path)
		w, resp := WrapResponseWriter(w)
		if r.Method == http.MethodOptions {
			return
		}
		h.ServeHTTP(w, r)
		log.Println(fmt.Sprintf("Status: %v", resp.StatusCode))
	}
}
