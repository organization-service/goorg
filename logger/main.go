package logger

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

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

var (
	DefaultWriter io.Writer = os.Stdout
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

func getClientIP(r *http.Request) string {
	clientIP := r.Header.Get("X-Forwarded-For")
	clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
	if clientIP == "" {
		clientIP = strings.TrimSpace(r.Header.Get("X-Real-IP"))
	}
	if clientIP != "" {
		return clientIP
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

func outPutLog(endTime time.Time, latency time.Duration, statusCode int, ipAddr, method, path string) {
	// 時間 ステータスコード レイテンシー IPアドレス HTTPメソッド パス
	fmt.Fprintf(
		DefaultWriter,
		"%v | %3d | %13v | %15s | %-7s | %s\n",
		endTime.Format("2006/01/02 - 15:04:05"),
		statusCode,
		latency,
		ipAddr,
		method,
		path,
	)
}

func Log(h interface{}) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		start := time.Now()
		path := r.URL.Path
		raw := r.URL.RawQuery
		method := r.Method
		rAddr := getClientIP(r)
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
		end := time.Now()
		latency := end.Sub(start)
		statusCode := resp.StatusCode
		if raw != "" {
			path = path + "?" + raw
		}
		outPutLog(end, latency, statusCode, rAddr, method, path)
	}
}
