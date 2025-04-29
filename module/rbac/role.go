package rbac

import (
	"net/http"

	"github.com/euiko/webapp/module/rbac/lib/role"
)

type (
	walker interface {
		Walk(walkFunc) error
	}

	walkerFunc func(walkFunc) error

	walkFunc func(method, route string, handler http.Handler) error

	permissionSupport interface {
		Permission() *role.Permission
	}
)

func (f walkerFunc) Walk(wf walkFunc) error {
	return f(wf)
}

func collectPermissions(walker walker) role.PermissionManager {
	// use 256 as the initial capacity of the slice
	// which mostly enough for most cases
	const initialCap = 256
	permissionsMap := make(map[int64]*role.Permission, initialCap)

	// walk and build the roles slice first
	walker.Walk(rolesCollectorWalkFunc(&permissionsMap))

	permissions := make([]*role.Permission, 0, len(permissionsMap))
	for _, permission := range permissionsMap {
		permissions = append(permissions, permission)
	}

	return role.NewPermissionManager(permissions)
}

func rolesCollectorWalkFunc(permissionsMap *map[int64]*role.Permission) walkFunc {
	return func(method, route string, handler http.Handler) error {
		permissionHandler, ok := handler.(permissionSupport)
		if !ok {
			return nil
		}

		permission := permissionHandler.Permission()
		if permission == nil {
			return nil
		}

		id := permission.ID()
		_, exists := (*permissionsMap)[id]
		// skip if already exists
		if exists {
			return nil
		}

		(*permissionsMap)[id] = permission
		return nil
	}
}

func collectEndpoints(walker walker, permissionsManager role.PermissionManager) []role.Endpoint {
	// use 1024 as the initial capacity of the slice
	// which mostly enough for most cases
	const initialCap = 1024
	endpoints := make([]role.Endpoint, 0, initialCap)

	// walk and build the endpoints slice
	walker.Walk(endpointsCollectorWalkFunc(&endpoints, permissionsManager))
	return endpoints
}

func endpointsCollectorWalkFunc(endpoints *[]role.Endpoint, permissionsManager role.PermissionManager) walkFunc {
	return func(method, route string, handler http.Handler) error {
		permissionHandler, ok := handler.(permissionSupport)
		if !ok {
			return nil
		}

		permission := permissionHandler.Permission()
		if permission == nil {
			return nil
		}

		if !permissionsManager.Has(permission) {
			return nil
		}

		endpointID := role.ID(method, route)
		*endpoints = append(*endpoints, role.Endpoint{
			ID:         endpointID,
			Method:     method,
			Path:       route,
			Permission: permission,
		})

		return nil
	}
}
