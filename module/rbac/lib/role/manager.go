package role

import "github.com/euiko/webapp/pkg/iter"

type (
	PermissionManager interface {
		Has(*Permission) bool
		HasAny(...*Permission) bool
		HasAll(...*Permission) bool
		HasID(int64) bool
		HasAnyIDs(...int64) bool
		HasAllIDs(...int64) bool
		All() []*Permission
		Map() map[int64]*Permission
	}

	permissionManager struct {
		permissions    []*Permission
		permissionsMap map[int64]*Permission
	}
)

func NewPermissionManager(permissions []*Permission) PermissionManager {
	permissionsMap := make(map[int64]*Permission, len(permissions))
	for _, p := range permissions {
		permissionsMap[p.ID()] = p
	}

	return &permissionManager{
		permissions:    permissions,
		permissionsMap: permissionsMap,
	}
}

// All implements PermissionManager.
func (m *permissionManager) All() []*Permission {
	return m.permissions
}

// Map implements PermissionManager.
func (m *permissionManager) Map() map[int64]*Permission {
	return m.permissionsMap
}

// Has implements PermissionManager.
func (m *permissionManager) Has(p *Permission) bool {
	return m.HasAny(p)
}

// HasAll implements PermissionManager.
func (m *permissionManager) HasAll(permissions ...*Permission) bool {
	permissionIDs := iter.Map(permissions, func(p *Permission) int64 {
		return p.ID()
	})
	return m.HasAllIDs(permissionIDs...)
}

// HasAny implements PermissionManager.
func (m *permissionManager) HasAny(permissions ...*Permission) bool {
	permissionIDs := iter.Map(permissions, func(p *Permission) int64 {
		return p.ID()
	})
	return m.HasAnyIDs(permissionIDs...)
}

// HasID implements PermissionManager.
func (m *permissionManager) HasID(id int64) bool {
	return m.HasAnyIDs(id)
}

// HasAnyIDs implements PermissionManager.
func (m *permissionManager) HasAnyIDs(permissionIDs ...int64) bool {
	for _, id := range permissionIDs {
		if _, ok := m.permissionsMap[id]; ok {
			return true
		}
	}

	return false
}

// HasAllIDs implements PermissionManager.
func (m *permissionManager) HasAllIDs(permissionIDs ...int64) bool {
	ok := iter.Reduce(permissionIDs, func(acc bool, id int64) bool {
		if _, ok := m.permissionsMap[id]; !ok {
			return false
		}

		return acc
	}, true)

	return ok
}
