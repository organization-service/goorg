package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/organization-service/goorg/v2/internal"
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
				ErrorHandler:  internal.OnError,
			})
		}
		if err := AuthValid.CheckJWT(rw, r); err != nil {
			log.Println(err)
		} else {
			internal.CallHandler(h, rw, r)
		}
	}
}
