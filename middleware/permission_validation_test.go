package middleware_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/organization-service/goorg/v2/auth"
	"github.com/organization-service/goorg/v2/logger"
	"github.com/organization-service/goorg/v2/middleware"
	"github.com/organization-service/goorg/v2/router"
	"github.com/stretchr/testify/assert"
)

func TestPermission(t *testing.T) {
	auth.AuthValid = jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(t *jwt.Token) (interface{}, error) {
			signBytes := []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDHW/yNyJrBa13q4dddbMmj+xP6
5K8QjCyOVJDU/glNVHgdCieHNeb26zzVH8a1HB2FrR7Ugcd6qRLW4xHg+nx3d8k7
5L+LJqJz/YsnDqP7dUmL2cYH1yKkbOFu5lRyf78tztZzlfZji+Dvk5RWmhNDj1iH
f2HhYlSbKH8Neje6cwIDAQAB
-----END PUBLIC KEY-----
`)
			return jwt.ParseRSAPublicKeyFromPEM(signBytes)
		},
		SigningMethod: jwt.SigningMethodRS256,
	})
	router := router.New()
	router.POST("/token", func(w http.ResponseWriter, r *http.Request) {
		signBytes := []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDHW/yNyJrBa13q4dddbMmj+xP65K8QjCyOVJDU/glNVHgdCieH
Neb26zzVH8a1HB2FrR7Ugcd6qRLW4xHg+nx3d8k75L+LJqJz/YsnDqP7dUmL2cYH
1yKkbOFu5lRyf78tztZzlfZji+Dvk5RWmhNDj1iHf2HhYlSbKH8Neje6cwIDAQAB
AoGACyhm1jioELNFqmPBfgUctATVdXFfKvntdHnfMUyWkLFtl5J9699kceVwni0N
Hg3YySaLVzF8TK0x1L0YvxLF7IBmCtvHSPXZjP+oSw0bJf9DTznqgbsXfcjPm3S/
p5cz5SPjOMzVkX2iiMDDjQYbm6coRwuydqg+4/ycwdl4CyECQQD7BaDI9X+3pHh6
JBps0OV8Fs7fVUCcTXh9Ql0l0Qo+lNodXzYVjz60erdxmqdptANnfqiu/B2BXD6A
WklcysVfAkEAy1AUszrHPJDoKokyx+rgwzoaIsymdIktCBdfna3RTXwZpApQwq4r
mNz6w5gg+ZxEJnXF+DazN/arsmvIbSLvbQJBANngCa1DQGZxz8wb1//I1NZ+qXI4
+cpwh3sZBeZT6TNmWWaTBEt0OHXH6b8l/9cEUswaqGixFR75pJodQ33R1bsCQQC4
hPJ0g7kE3+LFkAUVabcZl+bWhHPhugmzmTr/KRyXPTUsKuyG83m/33Z7A9uRjuBP
I35LuqFG/klvweCCJD21AkEA90+deyfHi/0wSzf2+KjCfj5YRTjSVGzPqmMGJ25a
hbbTf+WdLNehLWUasniVENojv/vk52ii7EoVHvFyrDMUsw==
-----END RSA PRIVATE KEY-----
`)
		signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
		if err != nil {
			logger.Log.Error(err)
		}
		token := jwt.New(jwt.SigningMethodRS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = "test"
		claims["sub"] = "abcdefghijklmn"
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		buf, _ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		mp := map[string]interface{}{}
		json.Unmarshal(buf, &mp)
		if mp != nil {
			if val, ok := mp["permission"]; ok {
				claims["permission"] = val
			}
		}
		tokenString, err := token.SignedString(signKey)
		if err != nil {
			logger.Log.Error(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(tokenString))
	})
	router.GET("/", auth.AuthMiddleware(middleware.PermissionCheck(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	}, "permission", "read:users", "create:users")))
	tests := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{
			name: "Permission Error",
			fn: func(t *testing.T) {
				req := httptest.NewRequest(http.MethodPost, "/token", nil)
				rw := httptest.NewRecorder()
				router.ServeHTTP(rw, req)
				assert.Equal(t, 200, rw.Result().StatusCode)
				buf, _ := ioutil.ReadAll(rw.Result().Body)
				defer rw.Result().Body.Close()
				assert.NotEqual(t, string(buf), "")
				req = httptest.NewRequest(http.MethodGet, "/", nil)
				rw = httptest.NewRecorder()
				req.Header.Add("Authorization", "Bearer "+string(buf))
				router.ServeHTTP(rw, req)
				assert.Equal(t, 401, rw.Result().StatusCode)
			},
		},
		{
			name: "Permission NG",
			fn: func(t *testing.T) {
				mp := map[string]interface{}{
					"permission": []string{"create:users"},
				}
				buf, _ := json.Marshal(&mp)
				req := httptest.NewRequest(http.MethodPost, "/token", bytes.NewBuffer(buf))
				rw := httptest.NewRecorder()
				router.ServeHTTP(rw, req)
				assert.Equal(t, 200, rw.Result().StatusCode)
				buf, _ = ioutil.ReadAll(rw.Result().Body)
				defer rw.Result().Body.Close()
				assert.NotEqual(t, string(buf), "")
				req = httptest.NewRequest(http.MethodGet, "/", nil)
				rw = httptest.NewRecorder()
				req.Header.Add("Authorization", "Bearer "+string(buf))
				router.ServeHTTP(rw, req)
				assert.Equal(t, 401, rw.Result().StatusCode)
				buf, _ = ioutil.ReadAll(rw.Result().Body)
				defer rw.Result().Body.Close()
			},
		},
		{
			name: "Permission OK",
			fn: func(t *testing.T) {
				mp := map[string]interface{}{
					"permission": []string{"read:users", "create:users"},
				}
				buf, _ := json.Marshal(&mp)
				req := httptest.NewRequest(http.MethodPost, "/token", bytes.NewBuffer(buf))
				rw := httptest.NewRecorder()
				router.ServeHTTP(rw, req)
				assert.Equal(t, 200, rw.Result().StatusCode)
				buf, _ = ioutil.ReadAll(rw.Result().Body)
				defer rw.Result().Body.Close()
				assert.NotEqual(t, string(buf), "")
				req = httptest.NewRequest(http.MethodGet, "/", nil)
				rw = httptest.NewRecorder()
				req.Header.Add("Authorization", "Bearer "+string(buf))
				router.ServeHTTP(rw, req)
				assert.Equal(t, 200, rw.Result().StatusCode)
				buf, _ = ioutil.ReadAll(rw.Result().Body)
				defer rw.Result().Body.Close()
				assert.Equal(t, string(buf), "OK")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.fn)
	}
}
