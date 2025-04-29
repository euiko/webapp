package lib

import (
	"github.com/euiko/webapp/module/rbac/lib/role"
)

var (
	PermissionManageRoles = role.Group("admin").NewPermission("manage-roles", "Manage roles")
)
