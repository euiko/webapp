package api

import (
	"time"

	"github.com/euiko/webapp/module/rbac/lib"
	"github.com/euiko/webapp/module/rbac/lib/role"
	"github.com/euiko/webapp/pkg/common/httpapi"
)

type (
	BaseRole struct {
		Name        string  `json:"name"`
		PrettyName  string  `json:"pretty_name"`
		Description string  `json:"description"`
		Permissions []int64 `json:"permissions"`
	}

	NewRole BaseRole

	UpdateRole struct {
		PrettyName  string  `json:"pretty_name"`
		Description string  `json:"description"`
		Permissions []int64 `json:"permissions"`
	}

	Role struct {
		BaseRole
		ID        int64     `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	ListAllRolesParams struct {
		httpapi.SearchParams
		httpapi.PaginationParams
	}

	ListAllRolesResponse struct {
		Items      []Role             `json:"items"`
		Pagination httpapi.Pagination `json:"pagination"`
	}
)

func ToRole(role role.Role) Role {
	return Role{
		BaseRole: BaseRole{
			Name:        role.Name,
			PrettyName:  role.PrettyName,
			Description: role.Description,
			Permissions: role.Permissions,
		},
		ID:        role.ID,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
	}
}

func ToListAllRolesResponse(pagy httpapi.PaginationParams, roles []role.Role, total int) ListAllRolesResponse {
	resp := ListAllRolesResponse{
		Items: make([]Role, len(roles)),
		Pagination: httpapi.Pagination{
			Page:     pagy.Page,
			PageSize: pagy.PageSize,
			Total:    total,
		},
	}

	for i, role := range roles {
		resp.Items[i] = ToRole(role)
	}

	return resp
}

func (p ListAllRolesParams) ToBase() lib.ListAllRolesParams {
	return lib.ListAllRolesParams{
		SearchParams:     p.SearchParams.ToBase(),
		PaginationParams: p.PaginationParams.ToBase(),
	}
}

func (p NewRole) ToBase() role.New {
	return role.New{
		Name:        p.Name,
		PrettyName:  p.PrettyName,
		Description: p.Description,
		Permissions: p.Permissions,
	}
}

func (p UpdateRole) ToBase() role.Update {
	return role.Update{
		PrettyName:  p.PrettyName,
		Description: p.Description,
		Permissions: p.Permissions,
	}
}
