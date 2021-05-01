package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/organization-service/goorg/v2/internal"
	"github.com/organization-service/goorg/v2/logger"
)

var getPemFn func(t *jwt.Token) (string, error)

func auth0Validation(c context.Context, t *jwt.Token, fn ...ValidaterFunc) (interface{}, error) {
	AUTH0_DOMAIN := os.Getenv("AUTH0_DOMAIN")
	if AUTH0_DOMAIN == "" {
		return "", errors.New("AUTH0_DOMAIN environment variables is empty")
	}
	url := fmt.Sprintf("https://%s/.well-known/jwks.json", AUTH0_DOMAIN)
	aud := os.Getenv("AUTH0_AUD")
	if aud != "" {
		if checkAud := t.Claims.(jwt.MapClaims).VerifyAudience(aud, false); !checkAud {
			return t, errors.New("Invalid audience.")
		}
	}
	iss := os.Getenv("AUTH0_ISS")
	if iss != "" {
		if checkIss := t.Claims.(jwt.MapClaims).VerifyIssuer(iss, false); !checkIss {
			return t, errors.New("Invalid issuer.")
		}
	}
	for _, f := range fn {
		if err := f(c, t); err != nil {
			return t, err
		}
	}
	cert, err := getPem(c, t, url)
	if err != nil {
		return t, err
	}
	result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
	return result, nil
}

func Auth0Validation(h interface{}) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		authValid := jwtmiddleware.New(jwtmiddleware.Options{
			ValidationKeyGetter: func(t *jwt.Token) (interface{}, error) {
				c := context.Background()
				return auth0Validation(c, t)
			},
			SigningMethod: jwt.SigningMethodRS256,
			ErrorHandler:  internal.OnError,
		})
		if err := authValid.CheckJWT(rw, r); err != nil {
			logger.Log.Error(err)
		} else {
			internal.CallHandler(h, rw, r)
		}
	}
}
