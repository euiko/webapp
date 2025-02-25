package cli

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/euiko/webapp/api"
	"github.com/euiko/webapp/db/sqldb"
	"github.com/euiko/webapp/settings"
	"github.com/spf13/cobra"

	// migration drivers
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlserver"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type (
	dbCmdConfig struct {
		openDb bool
	}

	dbCmdOption func(c *dbCmdConfig)
)

func Migration() api.Module {
	return api.NewModule(api.ModuleWithCLI(func(cmd *cobra.Command, s *settings.Settings) {
		cmd.AddCommand(migrationCmd(&s.DB.Sql))
	}))
}

func skipOpenDb() func(c *dbCmdConfig) {
	return func(c *dbCmdConfig) {
		c.openDb = false
	}
}

func migrationCmd(s *settings.SqlDatabase) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db",
		Short: "Database related commands",
	}
	cmd.AddCommand(newDbCmd(s, dbMigrateCmd))
	cmd.AddCommand(newDbCmd(s, dbRollbackCmd))
	cmd.AddCommand(newDbCmd(s, dbCreateMigrationCmd, skipOpenDb()))
	cmd.AddCommand(newDbCmd(s, dbLockCmd))
	cmd.AddCommand(newDbCmd(s, dbUnlockCmd))
	cmd.AddCommand(newDbCmd(s, dbStatusCmd))
	return cmd
}

func dbMigrateCmd(migrator *sqldb.Migrator, cmd *cobra.Command) {
	var steps int

	cmd.Use = "migrate"
	cmd.Short = "Migrate the database migrations"
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		migrated, err := (*migrator).Migrate(cmd.Context(), steps)
		if err == sqldb.ErrNoMigrations {
			fmt.Println("there are no migrations to be applied")
			return nil
		}

		if err != nil {
			fmt.Println("migration failed with error:", err)
		}

		fmt.Println("the following migrations were applied:", migrated)
		return err
	}
	cmd.Flags().IntVarP(&steps, "steps", "s", 0, "Number of migrations to apply")
}

func dbRollbackCmd(migrator *sqldb.Migrator, cmd *cobra.Command) {
	var steps int

	cmd.Use = "rollback"
	cmd.Short = "Rollback the database migration"
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		migrated, err := (*migrator).Rollback(cmd.Context(), steps)
		if err == sqldb.ErrNoMigrations {
			fmt.Println("there are no migrations to be applied")
			return nil
		}

		if err != nil {
			fmt.Println("rollback migration failed with error:", err)
		}

		fmt.Println("the following migrations were rolled back:", migrated)
		return err
	}

	cmd.Flags().IntVarP(&steps, "steps", "s", 0, "Number of migrations to revert")

}

func dbCreateMigrationCmd(migrator *sqldb.Migrator, cmd *cobra.Command) {
	var (
		dir string
	)
	cmd.Use = "create-migration [NAMES...]"
	cmd.Short = "Generate a new migration file"
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		// ensure the directory exists
		_ = os.MkdirAll(dir, 0755)

		// combine all the args to a single name with _ separator
		name := strings.Join(args, "_")
		fileNames, err := (*migrator).Create(cmd.Context(), name, dir)
		if err != nil {
			return err
		}

		for _, fileName := range fileNames {
			fmt.Printf("created migration %s (%s)\n", fileName, path.Join(dir, fileName))
		}

		return nil
	}

	cmd.Flags().StringVarP(&dir, "dir", "d", "./db/migrations", "Path to migration files")
}

func dbLockCmd(migrator *sqldb.Migrator, cmd *cobra.Command) {
	cmd.Use = "lock"
	cmd.Short = "Lock the database migrations"
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return (*migrator).Lock(cmd.Context())
	}
}

func dbUnlockCmd(migrator *sqldb.Migrator, cmd *cobra.Command) {
	cmd.Use = "unlock"
	cmd.Short = "Unlock the database migrations"
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return (*migrator).Unlock(cmd.Context())
	}
}

func dbStatusCmd(migrator *sqldb.Migrator, cmd *cobra.Command) {
	cmd.Use = "status"
	cmd.Short = "Show the database migrations status"
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		status, err := (*migrator).Status(cmd.Context())
		if err != nil {
			return err
		}

		fmt.Printf("migrations: %s\n", status.Migrations)
		fmt.Printf("unapplied migrations: %s\n", status.Unapplied)
		fmt.Printf("last applied migration: %s\n", status.LastApplied)
		return nil
	}
}

func newDbCmd(
	s *settings.SqlDatabase,
	f func(migrator *sqldb.Migrator, cmd *cobra.Command),
	opts ...dbCmdOption,
) *cobra.Command {
	var (
		migrator sqldb.Migrator
		config   = dbCmdConfig{
			openDb: true,
		}
		cmd = new(cobra.Command)
	)

	for _, opt := range opts {
		opt(&config)
	}

	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		var err error
		if config.openDb {
			if err = sqldb.Open(s); err != nil {
				return err
			}
			migrator, err = sqldb.NewMigrator(sqldb.ORM())
		} else {
			migrator, err = sqldb.NewMigrator(nil)
		}

		return err
	}

	f(&migrator, cmd)
	return cmd
}
