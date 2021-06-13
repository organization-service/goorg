package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/form3tech-oss/jwt-go"
	"github.com/organization-service/goorg/v2/auth"
	"github.com/organization-service/goorg/v2/internal"
)

type permissionMap map[string]bool

func (p *permissionMap) IsDenied() bool {
	for _, val := range *p {
		if !val {
			return true
		}
	}
	return false
}

func (p *permissionMap) copy() permissionMap {
	mp := permissionMap{}
	for key := range *p {
		mp[key] = false
	}
	return mp
}

func newPermissionMap(permissions ...string) permissionMap {
	mpPermission := permissionMap{}
	for _, permission := range permissions {
		mpPermission[permission] = false
	}
	return mpPermission
}

func PermissionCheck(h interface{}, permissionsClaim string, permissions ...string) http.HandlerFunc {
	if permissionsClaim == "" {
		panic(errors.New("Permissions claims is empty."))
	}
	if len(permissions) == 0 {
		panic(errors.New("Permissions is zero."))
	}
	originalPermission := newPermissionMap(permissions...)
	return func(w http.ResponseWriter, r *http.Request) {
		mpPermission := originalPermission.copy()
		ctx := r.Context()
		token := auth.FromContext(ctx)
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if claim, ok := claims[permissionsClaim]; ok {
				var ps []string
				switch values := claim.(type) {
				case string:
					// "read:test create:test"
					ps = strings.Split(values, " ")
				case []string:
					// ["read:test", "create:test"]
					ps = values
				case []interface{}:
					ps = make([]string, 0)
					for _, value := range values {
						if p, ok := value.(string); ok {
							ps = append(ps, p)
						}
					}
				}
				for _, p := range ps {
					if _, ok := mpPermission[p]; ok {
						mpPermission[p] = true
					}
				}
			}
		}
		if mpPermission.IsDenied() {
			internal.OnError(w, r, "Invalid token")
		} else {
			internal.CallHandler(h, w, r)
		}
	}
}
