package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/form3tech-oss/jwt-go"
	"github.com/organization-service/goorg/httpclient"
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

func getPem(ctx context.Context, t *jwt.Token, url string) (string, error) {
	cert := ""
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
		if t.Header["kid"] == val.Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + val.X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}
	if cert == "" {
		return cert, errors.New("Unable to find appropriate key.")
	}
	return cert, nil
}

func FromContext(c context.Context) *jwt.Token {
	return c.Value(AuthValid.Validation.Options.UserProperty).(*jwt.Token)
}
