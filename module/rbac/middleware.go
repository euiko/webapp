package rbac

import (
	"net/http"

	authlib "github.com/euiko/webapp/module/auth/lib"
	"github.com/euiko/webapp/module/rbac/lib"
)

func newMiddleware(m *Module) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := authlib.CurrentUser(r.Context())
			if !ok {
				unauthorized(w, r)
				return
			}

			userPermission, ok := user.(lib.User)
			if !ok {
				unauthorized(w, r)
				return
			}

			_ = userPermission.RoleName()

			// routeContext := chi.RouteContext(r.Context())
			// id := role.ID(routeContext.RouteMethod, routeContext.RoutePath)
			// endpoint, ok := m.endpointsMap[id]

			// // if there is no endpoint found in role, then we can skip the permission check
			// if !ok {
			// 	next.ServeHTTP(w, r)
			// 	return
			// }

		})
	}
}

func unauthorized(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("unauthorized"))
	w.WriteHeader(http.StatusUnauthorized)
}
