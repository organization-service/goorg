package router_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/organization-service/goorg/router"
	"github.com/stretchr/testify/assert"
)

func TestNewrelic(t *testing.T) {
	os.Setenv("APM_NAME", "newrelic")
	body := "newrelic"
	router := router.New()
	setRouter(router, body)
	tests := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{
			name: "get",
			fn: func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				rw := httptest.NewRecorder()
				router.ServeHTTP(rw, req)
				assert.Equal(t, 200, rw.Result().StatusCode)
				buf, _ := ioutil.ReadAll(rw.Result().Body)
				defer rw.Result().Body.Close()
				assert.Equal(t, []byte(body), buf)
			},
		},
		{
			name: "post",
			fn: func(t *testing.T) {
				req := httptest.NewRequest(http.MethodPost, "/", nil)
				rw := httptest.NewRecorder()
				router.ServeHTTP(rw, req)
				assert.Equal(t, 200, rw.Result().StatusCode)
				buf, _ := ioutil.ReadAll(rw.Result().Body)
				defer rw.Result().Body.Close()
				assert.Equal(t, []byte(body), buf)
			},
		},
		{
			name: "put",
			fn: func(t *testing.T) {
				req := httptest.NewRequest(http.MethodPut, "/", nil)
				rw := httptest.NewRecorder()
				router.ServeHTTP(rw, req)
				assert.Equal(t, 200, rw.Result().StatusCode)
				buf, _ := ioutil.ReadAll(rw.Result().Body)
				defer rw.Result().Body.Close()
				assert.Equal(t, []byte(body), buf)
			},
		},
		{
			name: "patch",
			fn: func(t *testing.T) {
				req := httptest.NewRequest(http.MethodPatch, "/", nil)
				rw := httptest.NewRecorder()
				router.ServeHTTP(rw, req)
				assert.Equal(t, 200, rw.Result().StatusCode)
				buf, _ := ioutil.ReadAll(rw.Result().Body)
				defer rw.Result().Body.Close()
				assert.Equal(t, []byte(body), buf)
			},
		},
		{
			name: "delete",
			fn: func(t *testing.T) {
				req := httptest.NewRequest(http.MethodDelete, "/", nil)
				rw := httptest.NewRecorder()
				router.ServeHTTP(rw, req)
				assert.Equal(t, 200, rw.Result().StatusCode)
				buf, _ := ioutil.ReadAll(rw.Result().Body)
				defer rw.Result().Body.Close()
				assert.Equal(t, []byte(body), buf)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.fn)
	}
}
