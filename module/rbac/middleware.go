package rbac

import (
	"net/http"

	authlib "github.com/euiko/webapp/module/auth/lib"
	"github.com/euiko/webapp/module/rbac/lib"
	"github.com/euiko/webapp/module/rbac/lib/role"
	"github.com/euiko/webapp/pkg/log"
	"github.com/go-chi/chi/v5"
)

func newMiddleware(m *Module) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// check if the target endpoint is protected by a role
			// if not, then we can skip the permission check
			routeContext := chi.RouteContext(r.Context())
			id := role.ID(routeContext.RouteMethod, routeContext.RoutePath)
			endpoint, ok := m.endpointsMap[id]
			if !ok {
				log.Info("skipping permission check for endpoint", log.WithField("method", routeContext.RouteMethod), log.WithField("path", routeContext.RoutePath))
				next.ServeHTTP(w, r) // skip permission check
				return
			}

			// check if the user is authenticated
			user, ok := authlib.CurrentUser(r.Context())
			if !ok {
				unauthorized(w, r)
				return
			}

			// ensure the user type is supporting a rbaclib.User
			userPermission, ok := user.(lib.User)
			if !ok {
				unauthorized(w, r)
				return
			}

			// obtain the permissions of the user
			// TODO: use better approach to avoid calling the database for each request
			roleName := userPermission.RoleName()
			userRole, err := m.GetRole(r.Context(), roleName)
			if err != nil {
				unauthorized(w, r)
				return
			}

			permissionManager := m.buildPermissionManager(userRole.Permissions...)
			ctx := lib.ContextWithPermissions(r.Context(), permissionManager)
			r = r.WithContext(ctx)

			// check if the user has the permission to access the endpoint
			if !permissionManager.Has(endpoint.Permission) {
				unauthorized(w, r)
				return
			}

			// call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

func unauthorized(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("unauthorized"))
	w.WriteHeader(http.StatusUnauthorized)
}
