package lib

import (
	"context"

	"github.com/euiko/webapp/module/rbac/lib/role"
)

type (
	Module interface {
		ListAllRoles(ctx context.Context, params ListAllRolesParams) ([]role.Role, int, error)
		GetRole(ctx context.Context, name string) (*role.Role, error)
		AddRole(ctx context.Context, r role.New) error
		RemoveRole(ctx context.Context, name string) error
		UpdateRole(ctx context.Context, name string, r role.Update) error
	}
)
