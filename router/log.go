package router

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/organization-service/goorg/v2/logger"
)

type handlerType int

const (
	httprouterHandle handlerType = iota
	handler
	handlerFunc
)

type requestInfo struct {
	StartTime     time.Time
	Path          string
	RawQuery      string
	Method        string
	RemoteAddress string
	EndTime       time.Time
}

func newRequestInfo(r *http.Request) *requestInfo {
	return &requestInfo{
		StartTime:     time.Now(),
		Path:          r.URL.Path,
		RawQuery:      r.URL.RawQuery,
		Method:        r.Method,
		RemoteAddress: getClientIP(r),
	}
}

func (i *requestInfo) RequestEnd() {
	i.EndTime = time.Now()
}
func (i *requestInfo) GetLatency() time.Duration {
	return i.EndTime.Sub(i.StartTime)
}

func (i *requestInfo) GetPath() string {
	if i.RawQuery != "" {
		return i.Path + "?" + i.RawQuery
	}
	return i.Path
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

func outPutLog(info *requestInfo, statusCode int) {
	// 時間 ステータスコード レイテンシー IPアドレス HTTPメソッド パス
	timeFormat := "2006/01/02 15:04:05"
	s := fmt.Sprintf(
		"code: %3d | time: %v - %v | latency: %13v | ip: %15s | method: %-7s | path: %s\n",
		statusCode,
		info.StartTime.Format(timeFormat),
		info.EndTime.Format(timeFormat),
		info.GetLatency(),
		info.RemoteAddress,
		info.Method,
		info.GetPath(),
	)
	logger.Log.Info(s)
}

func callHandler(h interface{}, w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	switch handler := h.(type) {
	case func(http.ResponseWriter, *http.Request, httprouter.Params):
		handler(w, r, p)
	case httprouter.Handle:
		handler(w, r, p)
	case http.HandlerFunc:
		*r = *r.WithContext(context.WithValue(r.Context(), httprouter.ParamsKey, p))
		handler.ServeHTTP(w, r)
	case http.Handler:
		*r = *r.WithContext(context.WithValue(r.Context(), httprouter.ParamsKey, p))
		handler.ServeHTTP(w, r)
	case func(http.ResponseWriter, *http.Request):
		*r = *r.WithContext(context.WithValue(r.Context(), httprouter.ParamsKey, p))
		handler(w, r)
	default:
		panic(errors.New("Not type handler"))
	}
}

func logHandler(h interface{}, typs ...handlerType) interface{} {
	typ := httprouterHandle
	if len(typs) >= 1 {
		typ = typs[0]
	}
	switch typ {
	case httprouterHandle:
		return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			info := newRequestInfo(r)
			w, resp := wrapResponseWriter(w)
			callHandler(h, w, r, p)
			info.RequestEnd()
			statusCode := resp.StatusCode
			outPutLog(info, statusCode)
		}
	case handler:
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			info := newRequestInfo(r)
			w, resp := wrapResponseWriter(w)
			callHandler(h, w, r, httprouter.Params{})
			info.RequestEnd()
			statusCode := resp.StatusCode
			outPutLog(info, statusCode)
		})
	case handlerFunc:
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			info := newRequestInfo(r)
			w, resp := wrapResponseWriter(w)
			callHandler(h, w, r, httprouter.Params{})
			info.RequestEnd()
			statusCode := resp.StatusCode
			outPutLog(info, statusCode)
		})
	default:
		panic(errors.New("Not type handler"))
	}
}
