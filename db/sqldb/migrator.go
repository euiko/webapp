package sqldb

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"regexp"
	"strings"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

type (
	Migrator interface {
		Migrate(ctx context.Context, steps ...int) ([]string, error)
		Rollback(ctx context.Context, steps ...int) ([]string, error)
		Lock(ctx context.Context) error
		Unlock(ctx context.Context) error
		Create(ctx context.Context, name, dir string) ([]string, error)
		Status(ctx context.Context) (*MigrationStatus, error)
		MarkApplied(ctx context.Context) error
	}

	MigratorOption func(c *migratorConfig)

	MigrationStatus struct {
		Migrations  fmt.Stringer
		Unapplied   fmt.Stringer
		LastApplied fmt.Stringer
	}

	migratorConfig struct {
		directories   []fs.FS
		markOnSuccess bool
	}

	bunMigrator struct {
		config   *migratorConfig
		migrator *migrate.Migrator
	}
)

var (
	ErrNoMigrations = errors.New("there are no migrations to be applied")

	globalDirectories = []fs.FS{}
	migrationRegex    = regexp.MustCompile(`[^a-zA-Z0-9_]+`)
)

func AddMigrationFS(fses ...fs.FS) {
	globalDirectories = append(globalDirectories, fses...)
}

func MigratorWithFS(fses ...fs.FS) MigratorOption {
	return func(c *migratorConfig) {
		c.directories = append(c.directories, fses...)
	}
}

func MigratorWithMarkOnSuccess(on bool) MigratorOption {
	return func(c *migratorConfig) {
		c.markOnSuccess = on
	}
}

func NewMigrator(db *bun.DB, opts ...MigratorOption) (Migrator, error) {
	config := newMigratorConfig()
	for _, opt := range opts {
		opt(config)
	}

	migrations := migrate.NewMigrations()
	directories := append(globalDirectories, config.directories...)
	for _, fsys := range directories {
		if err := migrations.Discover(fsys); err != nil {
			return nil, err
		}
	}

	m := bunMigrator{
		config: config,
		migrator: migrate.NewMigrator(db, migrations,
			migrate.WithMarkAppliedOnSuccess(config.markOnSuccess),
		),
	}

	return &m, nil
}

func (m *bunMigrator) Migrate(ctx context.Context, steps ...int) ([]string, error) {
	// due to the nature of the bun migrator implementation that groups
	// the migrations when being applied, we used a unique group id
	// for each migration, so we can advance by steps
	step := 0
	if len(steps) > 0 {
		step = steps[0]
	}

	// init migrations
	if err := m.migrator.Init(ctx); err != nil {
		return nil, err
	}

	applied, err := m.migrator.AppliedMigrations(ctx)
	if err != nil {
		return nil, err
	}
	lastGroupID := applied.LastGroupID()

	migrations, err := m.migrator.MigrationsWithStatus(ctx)
	if err != nil {
		return nil, err
	}
	migrations = migrations.Unapplied()

	if len(migrations) == 0 {
		return nil, ErrNoMigrations
	}

	migrated := make([]string, 0)
	for i := range migrations {
		// stop the migration when hit the step
		if step > 0 && i >= step {
			break
		}

		// generate new group id for every migrations
		groupID := lastGroupID + 1
		migration := &migrations[i]
		migration.GroupID = groupID

		if !m.config.markOnSuccess {
			if err := m.migrator.MarkApplied(ctx, migration); err != nil {
				return migrated, err
			}
		}

		if migration.Up != nil {
			if err := migration.Up(ctx, m.migrator.DB()); err != nil {
				return migrated, err
			}
		}

		if m.config.markOnSuccess {
			if err := m.migrator.MarkApplied(ctx, migration); err != nil {
				return migrated, err
			}
		}

		migrated = append(migrated, migration.Name)
		lastGroupID = groupID
	}

	return migrated, nil
}

func (m *bunMigrator) Rollback(ctx context.Context, steps ...int) ([]string, error) {
	// since every migration has unique group id, we can rollback by steps
	step := 0
	numMigrated := 0
	if len(steps) > 0 {
		step = steps[0]
	}

	migrations, err := m.migrator.MigrationsWithStatus(ctx)
	if err != nil {
		return nil, err
	}

	lastGroup := migrations.LastGroup()
	if len(lastGroup.Migrations) == 0 {
		return nil, ErrNoMigrations
	}

	migrated := make([]string, 0)
	for i := len(lastGroup.Migrations) - 1; i >= 0; i-- {
		// stop the migrations once hit the step
		if step > 0 && numMigrated >= step {
			break
		}

		migration := &lastGroup.Migrations[i]

		if !m.config.markOnSuccess {
			if err := m.migrator.MarkUnapplied(ctx, migration); err != nil {
				return migrated, err
			}
		}

		if migration.Down != nil {
			if err := migration.Down(ctx, m.migrator.DB()); err != nil {
				return migrated, err
			}
		}

		if m.config.markOnSuccess {
			if err := m.migrator.MarkUnapplied(ctx, migration); err != nil {
				return migrated, err
			}
		}

		// track the migrated migrations
		migrated = append(migrated, migration.Name)
		numMigrated++
	}

	return migrated, nil
}

func (m *bunMigrator) Lock(ctx context.Context) error {
	return m.migrator.Lock(ctx)
}

func (m *bunMigrator) Unlock(ctx context.Context) error {
	return m.migrator.Unlock(ctx)
}

func (m *bunMigrator) Create(ctx context.Context, name, dir string) ([]string, error) {
	// replace all space and dash with underscore
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")

	// remove all non-alphanumeric characters
	name = migrationRegex.ReplaceAllString(name, "")

	// convert to lower case
	name = strings.ToLower(name)

	// create new migrations and migrator just for this operation
	// to allow custom target directory
	migrations := migrate.NewMigrations(migrate.WithMigrationsDirectory(dir))
	migrator := migrate.NewMigrator(nil, migrations)

	// create migration files
	files, err := migrator.CreateSQLMigrations(ctx, name)
	if err != nil {
		return nil, err
	}

	// collect and return the generated file names
	fileNames := make([]string, len(files))
	for i, mf := range files {
		fileNames[i] = mf.Name
	}

	return fileNames, nil
}

func (m *bunMigrator) Status(ctx context.Context) (*MigrationStatus, error) {
	ms, err := m.migrator.MigrationsWithStatus(ctx)
	if err != nil {
		return nil, err
	}

	status := MigrationStatus{
		Migrations:  ms,
		Unapplied:   ms.Unapplied(),
		LastApplied: ms.LastGroup(),
	}
	return &status, nil
}

func (m *bunMigrator) MarkApplied(ctx context.Context) error {
	group, err := m.migrator.Migrate(ctx, migrate.WithNopMigration())
	if err != nil {
		return err
	}

	if group.IsZero() {
		return ErrNoMigrations
	}

	return nil
}

func newMigratorConfig() *migratorConfig {
	return &migratorConfig{
		directories:   []fs.FS{},
		markOnSuccess: true,
	}
}
