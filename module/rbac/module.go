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
		defaultRoles      []role.BaseRole
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

func WithDefaultRoles(roles ...role.BaseRole) ModuleOption {
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
		defaultRoles:      []role.BaseRole{},
		storeFactory: func(app core.App) Store {
			return NewOrmStore(sqldb.ORM())
		},
	}
}

func (m *Module) Init(ctx context.Context, s *settings.Settings) error {
	sqldb.AddMigrationFS(embededMigrationFS)

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

func (m *Module) GetPermissionByID(id int64) (*role.Permission, error) {
	if r, ok := m.permissionManager.Map()[id]; ok {
		return r, nil
	}

	return nil, errors.New("role not found")
}

func (m *Module) ListAllRoles(ctx context.Context, params lib.ListAllRolesParams) ([]role.Role, int, error) {
	if err := validator.Validate(params); err != nil {
		return nil, 0, err
	}

	return m.store.GetAll(ctx, params)
}

func (m *Module) GetRoleByName(ctx context.Context, name string) (*role.Role, error) {
	if name == "" {
		return nil, errors.New("role name is empty")
	}

	return m.store.GetRoleByName(ctx, name)
}

func (m *Module) buildEndpointsMap() {
	m.endpointsMap = make(map[int64]role.Endpoint, len(m.endpoints))
	for _, e := range m.endpoints {
		m.endpointsMap[e.ID] = e
	}
}

func (m *Module) ensureRoleExists(ctx context.Context, r role.BaseRole) error {
	// skip if already exists
	_, err := m.store.GetRoleByName(ctx, r.Name)
	if sqldb.IsNoRows(err) {
		return m.store.CreateRole(ctx, role.NewRole(r))
	}

	return err
}
