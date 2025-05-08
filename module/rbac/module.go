package rbac

import (
	"context"
	"embed"
	"errors"

	"github.com/euiko/webapp/core"
	"github.com/euiko/webapp/db/sqldb"
	"github.com/euiko/webapp/module/rbac/lib"
	"github.com/euiko/webapp/module/rbac/lib/role"
	"github.com/euiko/webapp/pkg/validator"
	"github.com/euiko/webapp/settings"
)

type (
	Module struct {
		app          core.App
		storeFactory StoreFactory

		permissionManager role.PermissionManager
		endpoints         []role.Endpoint
		endpointsMap      map[int64]role.Endpoint
		defaultRoles      []role.Base
		store             Store
	}

	ModuleOption func(*Module)

	StoreFactory func(app core.App) Store
)

var (
	//go:embed internal/migrations
	embededMigrationFS embed.FS
)

func ModuleFactory(options ...ModuleOption) core.ModuleFactory {
	return func(app core.App) core.Module {
		return NewModule(app, options...)
	}
}

func WithDefaultRoles(roles ...role.Base) ModuleOption {
	return func(m *Module) {
		m.defaultRoles = roles
	}
}

func WithStoreFactory(factory StoreFactory) ModuleOption {
	return func(m *Module) {
		m.storeFactory = factory
	}
}

func NewModule(app core.App, options ...ModuleOption) *Module {
	return &Module{
		app:               app,
		permissionManager: nil,
		endpoints:         make([]role.Endpoint, 0),
		endpointsMap:      map[int64]role.Endpoint{},
		defaultRoles:      []role.Base{},
		storeFactory: func(app core.App) Store {
			return NewOrmStore(sqldb.ORM())
		},
	}
}

func (m *Module) Init(ctx context.Context, s *settings.Settings) error {
	sqldb.AddMigrationFS(embededMigrationFS)
	m.app.AddMiddleware(newMiddleware(m))

	return nil
}

func (m *Module) Close() error {
	return nil
}

func (m *Module) BeforeStart(ctx context.Context) error {
	m.store = m.storeFactory(m.app)

	for _, role := range m.defaultRoles {
		if err := m.ensureRoleExists(ctx, role); err != nil {
			return err
		}
	}

	return nil
}

func (m *Module) ListAllRoles(ctx context.Context, params lib.ListAllRolesParams) ([]role.Role, int, error) {
	if err := validator.Validate(params); err != nil {
		return nil, 0, err
	}

	return m.store.GetAll(ctx, params)
}

func (m *Module) GetRole(ctx context.Context, name string) (*role.Role, error) {
	if name == "" {
		return nil, errors.New("role name is empty")
	}

	return m.store.Get(ctx, name)
}

func (m *Module) AddRole(ctx context.Context, r role.New) error {
	if err := validator.Validate(r); err != nil {
		return err
	}

	// ensure all role permissions are valid
	if !m.permissionManager.HasAllIDs(r.Permissions...) {
		return errors.New("invalid role permissions")
	}

	return m.store.Create(ctx, r)
}

func (m *Module) RemoveRole(ctx context.Context, name string) error {
	if name == "" {
		return errors.New("role name is empty")
	}

	return m.store.Delete(ctx, name)
}

func (m *Module) UpdateRole(ctx context.Context, name string, r role.Update) error {
	if name == "" {
		return errors.New("role name is empty")
	}

	if err := validator.Validate(r); err != nil {
		return err
	}

	// ensure all role permissions are valid
	if !m.permissionManager.HasAllIDs(r.Permissions...) {
		return errors.New("invalid role permissions")
	}

	return m.store.Update(ctx, name, r)
}

func (m *Module) buildPermissionManager(permissionIDs ...int64) role.PermissionManager {
	var (
		permissions   = make([]*role.Permission, 0, len(permissionIDs))
		permissionMap = m.permissionManager.Map()
	)

	for _, id := range permissionIDs {
		if p, ok := permissionMap[id]; ok {
			permissions = append(permissions, p)
		}
	}

	return role.NewPermissionManager(permissions)
}

func (m *Module) buildEndpointsMap() {
	m.endpointsMap = make(map[int64]role.Endpoint, len(m.endpoints))
	for _, e := range m.endpoints {
		m.endpointsMap[e.ID] = e
	}
}

func (m *Module) ensureRoleExists(ctx context.Context, r role.Base) error {
	// skip if already exists
	_, err := m.store.Get(ctx, r.Name)
	if sqldb.IsNoRows(err) {
		return m.store.Create(ctx, role.New(r))
	}

	return err
}
