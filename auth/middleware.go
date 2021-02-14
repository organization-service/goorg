package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/julienschmidt/httprouter"
)

var AuthValid *jwtmiddleware.JWTMiddleware

func AuthMiddleware(h interface{}) http.HandlerFunc {
	IDP_NAME := os.Getenv("IDP_NAME")
	return func(rw http.ResponseWriter, r *http.Request) {
		if nil == AuthValid {
			AuthValid = jwtmiddleware.New(jwtmiddleware.Options{
				ValidationKeyGetter: func(t *jwt.Token) (interface{}, error) {
					switch IDP_NAME {
					case "auth0":
						c := context.Background()
						return auth0Validation(c, t)
					default:
						return nil, fmt.Errorf("No support IDP_NAME: %s", IDP_NAME)
					}
				},
				SigningMethod: jwt.SigningMethodRS256,
			})
		}
		if err := AuthValid.CheckJWT(rw, r); err != nil {
			log.Println(err)
		} else {
			switch handler := h.(type) {
			case func(http.ResponseWriter, *http.Request, httprouter.Params):
				p := httprouter.ParamsFromContext(r.Context())
				handler(rw, r, p)
			case httprouter.Handle:
				p := httprouter.ParamsFromContext(r.Context())
				handler(rw, r, p)
			case http.HandlerFunc:
				handler.ServeHTTP(rw, r)
			case func(http.ResponseWriter, *http.Request):
				handler(rw, r)
			}
		}
	}
}
