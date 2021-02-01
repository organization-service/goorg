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

type (
	AuthCheck struct {
		Validation      *jwtmiddleware.JWTMiddleware
		ValidationFuncs []ValidaterFunc
		GetPem          GetPemFun
	}
)

var AuthValid = &AuthCheck{
	GetPem: getPem,
}

func AuthMiddleware(h interface{}) http.HandlerFunc {
	IDP_NAME := os.Getenv("IDP_NAME")
	if nil == AuthValid.ValidationFuncs {
		AuthValid.ValidationFuncs = append(AuthValid.ValidationFuncs, func(c context.Context, token *jwt.Token) error {
			return nil
		})
	}
	return func(rw http.ResponseWriter, r *http.Request) {
		c := r.Context()
		if nil == AuthValid.Validation {
			AuthValid.Validation = jwtmiddleware.New(jwtmiddleware.Options{
				ValidationKeyGetter: func(t *jwt.Token) (interface{}, error) {
					switch IDP_NAME {
					case "auth0":
						return auth0Validation(c, t, AuthValid.ValidationFuncs...)
					default:
						return nil, fmt.Errorf("No support IDP_NAME: %s", IDP_NAME)
					}
				},
				SigningMethod: jwt.SigningMethodRS256,
			})
		}
		if err := AuthValid.Validation.CheckJWT(rw, r); err != nil {
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
