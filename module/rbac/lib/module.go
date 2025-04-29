package lib

import (
	"context"

	"github.com/euiko/webapp/module/rbac/lib/role"
)

type (
	Module interface {
		GetPermissionByID(id int64) (*role.Permission, error)
		GetRoleByName(ctx context.Context, name string) (*role.Role, error)
		ListAllRoles(ctx context.Context, params ListAllRolesParams) ([]role.Role, int, error)
	}
)
