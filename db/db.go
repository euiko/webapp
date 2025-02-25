package db

import (
	"github.com/euiko/webapp/db/sqldb"
	"github.com/euiko/webapp/settings"
)

func Init(s *settings.Database) error {
	if err := initSqlDb(&s.Sql, s.Extra.Sql); err != nil {
		return err
	}

	return nil
}

func Close() error {
	var err error
	// keep try closes another db connection even when error
	// and return the last errors
	if e := sqldb.Close(); e != nil {
		err = e
	}

	return err
}

func initSqlDb(s *settings.SqlDatabase, extra map[string]settings.SqlDatabase) error {
	if err := sqldb.Open(s); err != nil {
		return err
	}

	for name, s := range extra {
		if err := sqldb.Open(&s, name); err != nil {
			return err
		}
	}

	return nil
}
