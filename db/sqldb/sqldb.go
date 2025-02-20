package sqldb

import (
	"database/sql"
	"errors"

	"github.com/euiko/webapp/pkg/log"
	"github.com/euiko/webapp/settings"
)

const defaultDbName = "default"

var (
	// for holds all the database connections
	instances    = make(map[string]*sql.DB)
	errNotOpened = errors.New("database not opened")
)

// DB returns the selected database connection by its name, default to default connection
func DB(names ...string) *sql.DB {
	name := defaultDbName
	if len(names) > 0 {
		name = names[0]
	}

	// ensure the database is opened
	if _, ok := instances[name]; !ok {
		log.Error(errNotOpened.Error(),
			log.WithField("name", name),
			log.WithError(errNotOpened),
		)
		// exit early
		panic(errNotOpened)
	}

	return instances[name]
}

func Open(s settings.SqlDatabase, names ...string) error {
	var (
		name = defaultDbName
	)

	// skip if not enabled
	if !s.Enabled {
		return nil
	}

	// use supplied name if any
	if len(names) > 0 {
		name = names[0]
	}

	// ensure the database is not already opened
	if _, ok := instances[name]; ok {
		return errors.New("database already opened")
	}

	// TODO: add gorm configurations
	db, err := sql.Open(s.Driver, s.Uri)
	if err != nil {
		return err
	}

	db.SetMaxIdleConns(s.MaxIdleConns)
	db.SetMaxOpenConns(s.MaxOpenConns)
	db.SetConnMaxLifetime(s.ConnMaxLifetime)

	if s.UseORM {
		if err := initORM(name, &s, db); err != nil {
			return err
		}
	}

	instances[name] = db
	return nil
}

func Close() error {
	var err error
	// keep try closes another db connection even when error
	// and return the last errors
	for name, db := range instances {
		if e := db.Close(); e != nil {
			log.Error("error when closing database",
				log.WithField("name", name),
				log.WithError(err))
			err = e // return the last error as errors
		}
	}

	return err
}
