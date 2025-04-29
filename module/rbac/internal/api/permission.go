package api

import "github.com/euiko/webapp/module/rbac/lib/role"

type (
	Permission struct {
		Group      string `json:"group"`
		Name       string `json:"name"`
		PrettyName string `json:"pretty_name"`
		ID         int64  `json:"id"`
	}
)

func ToPermission(p *role.Permission) Permission {
	return Permission{
		Group:      p.Group,
		Name:       p.Name,
		PrettyName: p.PrettyName,
		ID:         p.ID(),
	}
}

func ToPermissions(permissions ...*role.Permission) []Permission {
	perms := make([]Permission, len(permissions))
	for i, p := range permissions {
		perms[i] = ToPermission(p)
	}
	return perms
}
