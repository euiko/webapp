package role

type (
	PermissionManager interface {
		Has(*Permission) bool
		HasAny(...*Permission) bool
		HasAll(...*Permission) bool
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
	for _, p := range permissions {
		if _, ok := m.permissionsMap[p.ID()]; !ok {
			return false
		}
	}

	return true
}

// HasAny implements PermissionManager.
func (m *permissionManager) HasAny(permissions ...*Permission) bool {
	for _, p := range permissions {
		if _, ok := m.permissionsMap[p.ID()]; ok {
			return true
		}
	}

	return false
}
