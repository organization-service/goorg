package auth_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/form3tech-oss/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/organization-service/goorg/v2/auth"
	"github.com/organization-service/goorg/v2/router"
)

func TestAuth0(t *testing.T) {
	route := router.New()
	route.GET("/id/:id", auth.AuthMiddleware(func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		rw.Header().Set("Content-Type", "application/json")
		token := auth.FromContext(r.Context())
		t := token.Claims.(jwt.MapClaims)
		rw.Write([]byte(t["sub"].(string)))
	}))
	tests := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{
			name: "get1",
			fn: func(t *testing.T) {
				header := http.Header{}
				header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ID_TOKEN")))
				req := httptest.NewRequest(http.MethodGet, "/id/123", nil)
				req.Header = header
				rw := httptest.NewRecorder()
				route.ServeHTTP(rw, req)

			},
		},
		{
			name: "get2",
			fn: func(t *testing.T) {
				header := http.Header{}
				header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ID_TOKEN")))
				req := httptest.NewRequest(http.MethodGet, "/id/123", nil)
				req.Header = header
				rw := httptest.NewRecorder()
				route.ServeHTTP(rw, req)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.fn)
	}
}
