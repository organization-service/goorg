package middleware

import (
	"net/http"

	"github.com/form3tech-oss/jwt-go"
	"github.com/organization-service/goorg/auth"
	"github.com/organization-service/goorg/internal"
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

func PermissionCheck(h interface{}, permissionClaim string, permissions ...string) http.HandlerFunc {
	if permissionClaim == "" {
		return func(w http.ResponseWriter, r *http.Request) { internal.CallHandler(h, w, r) }
	}
	originalPermission := newPermissionMap(permissions...)
	return func(w http.ResponseWriter, r *http.Request) {
		mpPermission := originalPermission.copy()
		ctx := r.Context()
		token := auth.FromContext(ctx)
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if claim, ok := claims[permissionClaim]; ok {
				var ps []string
				switch values := claim.(type) {
				case []string:
					ps = make([]string, len(values))
					for idx, val := range values {
						ps[idx] = val
					}
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
