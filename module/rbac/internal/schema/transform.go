package schema

import "github.com/euiko/webapp/module/rbac/lib/role"

func (s Role) ToBase() role.Role {
	return role.Role{
		BaseRole: role.BaseRole{
			Name:        s.Name,
			PrettyName:  s.PrettyName,
			Description: s.Description,
			Permissions: s.Permissions,
		},
		ID:        s.ID,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}
