package rbac

import (
	"context"

	"github.com/euiko/webapp/db/sqldb"
	"github.com/euiko/webapp/module/rbac/internal/schema"
	"github.com/euiko/webapp/module/rbac/lib"
	"github.com/euiko/webapp/module/rbac/lib/role"
)

type (
	Store interface {
		GetAll(ctx context.Context, params lib.ListAllRolesParams) ([]role.Role, int, error)
		Get(ctx context.Context, name string) (*role.Role, error)
		Create(ctx context.Context, role role.New) error
		Delete(ctx context.Context, name string) error
		Update(ctx context.Context, name string, role role.Update) error
	}

	ormStore struct {
		db sqldb.OrmDB
	}
)

func NewOrmStore(db sqldb.OrmDB) Store {
	return &ormStore{
		db: db,
	}
}

// GetAll implements Store.
func (s *ormStore) GetAll(ctx context.Context, params lib.ListAllRolesParams) ([]role.Role, int, error) {
	var roles []schema.Role
	query := s.db.NewSelect().
		Model(&roles)

	if params.Keyword != "" {
		query = query.Where("name LIKE ?", "%"+params.Keyword+"%")
	}

	count, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	err = query.
		Limit(params.Limit).
		Offset(params.Offset).
		Scan(ctx)
	if err != nil {
		return nil, 0, err
	}

	result := make([]role.Role, len(roles))
	for i, role := range roles {
		result[i] = role.ToBase()
	}

	return result, count, nil
}

// GetRoleByName implements Store.
func (s *ormStore) Get(ctx context.Context, name string) (*role.Role, error) {
	var role schema.Role
	query := s.db.NewSelect().
		Model(&role).
		Where("name = ?", name).
		Limit(1)

	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}

	r := role.ToBase()
	return &r, nil
}

// Create implements Store.
func (s *ormStore) Create(ctx context.Context, role role.New) error {
	newRole := schema.Role{
		Name:        role.Name,
		PrettyName:  role.PrettyName,
		Description: role.Description,
		Permissions: role.Permissions,
	}

	_, err := s.db.NewInsert().
		Model(&newRole).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Update implements Store.
func (s *ormStore) Update(ctx context.Context, name string, role role.Update) error {
	updatedRole := schema.Role{
		PrettyName:  role.PrettyName,
		Description: role.Description,
		Permissions: role.Permissions,
	}

	_, err := s.db.NewUpdate().
		Model(&updatedRole).
		Where("name = ?", name).
		Exec(ctx)

	return err
}

// Delete implements Store.
func (s *ormStore) Delete(ctx context.Context, name string) error {
	_, err := s.db.NewDelete().
		Model((*schema.Role)(nil)).
		Where("name = ?", name).
		Exec(ctx)
	return err
}
