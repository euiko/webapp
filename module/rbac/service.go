package rbac

import (
	"errors"
	"net/http"

	"github.com/euiko/webapp/core"
	authlib "github.com/euiko/webapp/module/auth/lib"
	api "github.com/euiko/webapp/module/rbac/internal/api"
	"github.com/euiko/webapp/module/rbac/lib"
	"github.com/euiko/webapp/module/rbac/lib/role"
	"github.com/euiko/webapp/pkg/common/httpapi"
	"github.com/euiko/webapp/pkg/helper"
	"github.com/euiko/webapp/pkg/log"
	"github.com/go-chi/chi/v5"
)

func (m *Module) PostRoute(r core.Router) {
	walker := chiWalker(r)
	m.permissionManager = collectPermissions(walker)
	m.endpoints = collectEndpoints(walker, m.permissionManager)
	m.buildEndpointsMap()

	// when default roles is not defined then we will set the default roles
	// to be admin with all permissions
	if len(m.defaultRoles) == 0 {
		m.defaultRoles = append(m.defaultRoles, role.Base{
			Name:        "admin",
			PrettyName:  "Administrator",
			Description: "Role with all permissions",
			Permissions: role.PermissionsToIDs(m.permissionManager.All()...),
		})
	}
}

func (m *Module) APIRoute(r core.Router) {
	r.Group(func(r core.Router) {
		r.Use(authlib.AuthRequiredMiddleware(m.app))
		r.Get("/permissions", m.listAllPermissionsHandler)
		r.Method("GET", "/users/me/role", m.getRoleHandler(m.userFromSession))

		// endpoint that requires manage roles permission
		r.Get("/roles", m.listAllRolesHandler)
		r.Method("POST", "/roles", role.Handler(lib.PermissionManageRoles, http.HandlerFunc(m.addRoleHandler)))
		r.Method("DELETE", "/roles/{name}", role.Handler(lib.PermissionManageRoles, http.HandlerFunc(m.removeRoleHandler)))
		r.Method("PUT", "/roles/{name}", role.Handler(lib.PermissionManageRoles, http.HandlerFunc(m.updateRoleHandler)))
		r.Method("GET", "/users/{id}/role", role.Handler(lib.PermissionManageRoles, m.getRoleHandler(m.userFromIDParams)))
	})
}

func (m *Module) listAllPermissionsHandler(w http.ResponseWriter, r *http.Request) {
	permissions := m.permissionManager.All()
	helper.WriteResponse(w, api.ToPermissions(permissions...))
}

func (m *Module) listAllRolesHandler(w http.ResponseWriter, r *http.Request) {
	params := api.ListAllRolesParams{
		PaginationParams: httpapi.PaginationParams{
			Page:     1,
			PageSize: 10,
		},
	}
	if err := helper.DecodeRequest(r, &params); err != nil {
		helper.WriteResponse(w, err)
		return
	}

	roles, total, err := m.ListAllRoles(r.Context(), params.ToBase())
	if err != nil {
		helper.WriteResponse(w, err)
		return
	}

	response := api.ToListAllRolesResponse(params.PaginationParams, roles, total)
	helper.WriteResponse(w, response)
}

func (m *Module) getRoleHandler(userProvider func(*http.Request) authlib.User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := userProvider(r)
		if user == nil {
			helper.WriteResponse(w, errors.New("invalid user"), helper.ResponseWithStatus(http.StatusBadRequest))
			return
		}

		// get the role of the user
		roleUser, ok := user.(lib.User)
		if !ok {
			helper.WriteResponse(w, errors.New("user is not a role user"), helper.ResponseWithStatus(http.StatusInternalServerError))
			return
		}

		roleName := roleUser.RoleName()
		role, err := m.GetRole(r.Context(), roleName)
		if err != nil {
			helper.WriteResponse(w, err)
			return
		}

		response := api.ToRole(*role)
		helper.WriteResponse(w, response)
	})
}

func (m *Module) addRoleHandler(w http.ResponseWriter, r *http.Request) {
	var payload api.NewRole
	if err := helper.DecodeRequestBody(r, &payload); err != nil {
		helper.WriteResponse(w, err)
		return
	}

	if err := m.AddRole(r.Context(), payload.ToBase()); err != nil {
		helper.WriteResponse(w, err)
		return
	}

	helper.WriteResponse(w, "created")
}

func (m *Module) removeRoleHandler(w http.ResponseWriter, r *http.Request) {
	roleName := chi.URLParam(r, "name")
	if err := m.RemoveRole(r.Context(), roleName); err != nil {
		helper.WriteResponse(w, err)
		return
	}

	helper.WriteResponse(w, "deleted")
}

func (m *Module) updateRoleHandler(w http.ResponseWriter, r *http.Request) {
	var (
		payload  api.UpdateRole
		roleName = chi.URLParam(r, "name")
	)

	if err := helper.DecodeRequestBody(r, &payload); err != nil {
		helper.WriteResponse(w, err)
		return
	}

	if err := m.UpdateRole(r.Context(), roleName, payload.ToBase()); err != nil {
		helper.WriteResponse(w, err)
		return
	}

	helper.WriteResponse(w, "updated")
}

func (m *Module) userFromSession(r *http.Request) authlib.User {
	user, ok := authlib.CurrentUser(r.Context())
	if !ok {
		return nil
	}

	return user
}

func (m *Module) userFromIDParams(r *http.Request) authlib.User {
	userID := chi.URLParam(r, "id")
	if userID == "" {
		return nil
	}

	authModule := core.MustGetModule[authlib.Module](m.app)
	userLoader := authModule.UserLoader()

	user, err := userLoader.UserById(r.Context(), userID)
	if err != nil {
		log.Error("failed to load user", log.WithError(err))
		return nil
	}

	return user
}

func chiWalker(r chi.Routes) walker {
	return walkerFunc(func(wf walkFunc) error {
		return chi.Walk(r, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
			return wf(method, route, handler)
		})
	})
}
