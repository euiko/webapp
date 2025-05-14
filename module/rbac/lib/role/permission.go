package role

import (
	"github.com/ettle/strcase"
)

type (
	Endpoint struct {
		ID         int64
		Method     string
		Path       string
		Permission *Permission
	}

	Permission struct {
		Group      string
		Name       string
		PrettyName string

		id int64
	}

	Group string
)

func PermissionsToIDs(permissions ...*Permission) []int64 {
	ids := make([]int64, len(permissions))
	for i, p := range permissions {
		ids[i] = p.ID()
	}
	return ids
}

func NewPermission(group, name string, prettyNames ...string) *Permission {
	return &Permission{
		Group:      group,
		Name:       name,
		PrettyName: strcase.ToCase(name, strcase.Original, ' '),
		id:         ID(group, name),
	}
}

func (g Group) NewPermission(name string, prettyNames ...string) *Permission {
	return NewPermission(string(g), name, prettyNames...)
}

func (p *Permission) ID() int64 {
	return p.id
}
