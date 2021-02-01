package auth_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/form3tech-oss/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/organization-service/goorg/auth"
	"github.com/organization-service/goorg/httpclient"
	"github.com/organization-service/goorg/router"
)

func TestAuth0(t *testing.T) {
	route := router.New()
	route.GET("/id/:id", auth.AuthMiddleware(func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		token := auth.FromContext(r.Context())
		t := token.Claims.(jwt.MapClaims)
		log.Println(p.ByName("id"))
		rw.Write([]byte(t["sub"].(string)))
	}))
	c := context.Background()
	header := http.Header{}
	header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ID_TOKEN")))
	client := httpclient.NewClient(c, http.MethodGet, "/id/123", nil, header)
	req := client.GetRequest()
	rw := httptest.NewRecorder()
	route.ServeHTTP(rw, req)
}
