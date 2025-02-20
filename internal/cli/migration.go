package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/euiko/webapp/api"
	"github.com/euiko/webapp/settings"
	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"

	// migration drivers
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlserver"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	migrationPath string
	regex         = regexp.MustCompile(`[^a-zA-Z0-9_]+`)
)

func Migration(s *settings.Settings) api.Module {
	return api.NewModule(api.ModuleWithCLI(func(cmd *cobra.Command) {
		cmd.AddCommand(migrationCmd(s))
	}))
}

func migrationCmd(s *settings.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migration",
		Short: "Migrate the database",
	}
	cmd.PersistentFlags().StringVarP(&migrationPath, "path", "p", "./db/migrations", "Path to migration files")
	cmd.AddCommand(migrationUpCmd(s))
	cmd.AddCommand(migrationDownCmd(s))
	cmd.AddCommand(migrationNewCmd())
	return cmd
}

func migrationUpCmd(s *settings.Settings) *cobra.Command {
	var (
		steps int
		force int
	)

	cmd := &cobra.Command{
		Use:   "up",
		Short: "Migrate the database up",
		RunE: func(cmd *cobra.Command, args []string) error {
			source := fmt.Sprintf("file://%s", migrationPath)
			m, err := migrate.New(source, s.DB.Sql.Uri)
			if err != nil {
				return err
			}

			// force migration
			if force > 0 {
				if err := m.Force(force); err != nil {
					return err
				}
			}

			if steps > 0 {
				err = m.Steps(steps)
			} else {
				err = m.Up()
			}

			if err != migrate.ErrNoChange {
				return err
			}

			return nil
		},
	}
	cmd.Flags().IntVarP(&steps, "steps", "s", 0, "Number of migrations to apply")
	cmd.Flags().IntVarP(&force, "force", "f", 0, "Force specific version to apply")
	return cmd
}

func migrationDownCmd(s *settings.Settings) *cobra.Command {
	var steps int

	cmd := &cobra.Command{
		Use:   "down",
		Short: "Migrate the database down",
		RunE: func(cmd *cobra.Command, args []string) error {
			source := fmt.Sprintf("file://%s", migrationPath)
			m, err := migrate.New(source, s.DB.Sql.Uri)
			if err != nil {
				return err
			}

			if steps > 0 {
				// negate the steps to revert migrations
				return m.Steps(-1 * steps)
			} else {
				return m.Down()
			}
		},
	}
	cmd.Flags().IntVarP(&steps, "steps", "s", 0, "Number of migrations to revert")
	return cmd
}

func migrationNewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new [NAME]",
		Short: "Generate a new migration file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			return createMigrationFile(migrationPath, name)
		},
	}

	return cmd
}

func createMigrationFile(dir string, name string) error {
	// replace all space and dash with underscore
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")

	// remove all non-alphanumeric characters
	name = regex.ReplaceAllString(name, "")

	// convert to lower case
	name = strings.ToLower(name)

	if name == "" {
		return errors.New("you need to specify the migration name (alphanumeric with lowercase)")
	}

	// ensure directory exists
	if err := os.MkdirAll(dir, os.ModeDir); err != nil {
		return err
	}

	// add version
	now := time.Now()
	version := now.Unix()
	title := fmt.Sprintf("%d_%s", version, name)
	user := os.Getenv("USER")

	if user == "" {
		user = "-"
	}

	for _, direction := range []string{"up", "down"} {
		filename := fmt.Sprintf("%s.%s.sql", title, direction)
		file, err := os.Create(filepath.Join(dir, filename))
		if err != nil {
			return err
		}
		defer file.Close()

		fmt.Fprintf(file, "-- %s\n", filename)
		fmt.Fprintf(file, "-- Created at %s\n", now.Format(time.RFC3339))
		fmt.Fprintf(file, "-- By %s\n", user)
		fmt.Fprintf(file, "\n")
	}

	return nil
}
