package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/organization-service/goorg/v2/cache"
	"github.com/organization-service/goorg/v2/httpclient"
)

type (
	GetPemFun     func(ctx context.Context, t *jwt.Token, url string) (string, error)
	ValidaterFunc func(c context.Context, token *jwt.Token) error
	Jwks          struct {
		Keys []JSONWebKeys `json:"keys"`
	}
	JSONWebKeys struct {
		Kty string   `json:"kty"`
		Kid string   `json:"kid"`
		Use string   `json:"use"`
		N   string   `json:"n"`
		E   string   `json:"e"`
		X5c []string `json:"x5c"`
	}
)

var cachePem = cache.New()
var userProperty = "auth-user"

func getPem(ctx context.Context, t *jwt.Token, url string) (string, error) {
	kid := t.Header["kid"].(string)
	cert := ""
	cacheCert := cachePem.Get(kid)
	if cacheCert != nil {
		if val, ok := cacheCert.(string); ok {
			cert = "-----BEGIN CERTIFICATE-----\n" + val + "\n-----END CERTIFICATE-----"
			return cert, nil
		}
	}

	client := httpclient.NewClient(ctx, http.MethodGet, url, nil, http.Header{})
	resp, err := client.Do()
	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()
	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)
	if err != nil {
		return cert, err
	}
	for _, val := range jwks.Keys {
		if cachePem.Get(val.Kid) == nil {
			cachePem.Put(val.Kid, val.X5c[0], time.Now().Add(5*time.Minute).UnixNano())
		}
		if t.Header["kid"].(string) == val.Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + val.X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}
	if cert == "" {
		return cert, errors.New("Unable to find appropriate key.")
	}
	return cert, nil
}

func FromContext(c context.Context) *jwt.Token {
	return c.Value(userProperty).(*jwt.Token)
}
