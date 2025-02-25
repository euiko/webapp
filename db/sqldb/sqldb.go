package sqldb

import (
	"database/sql"
	"errors"

	"github.com/euiko/webapp/pkg/log"
	"github.com/euiko/webapp/settings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/microsoft/go-mssqldb"
)

const defaultDbName = "default"

var (
	ErrAlreadyOpened = errors.New("database already opened")
	ErrNotOpened     = errors.New("database not opened")

	// for holds all the database connections
	instances = make(map[string]*sql.DB)
)

// DB returns the selected database connection by its name, default to default connection
func DB(names ...string) *sql.DB {
	name := defaultDbName
	if len(names) > 0 {
		name = names[0]
	}

	// ensure the database is opened
	if _, ok := instances[name]; !ok {
		log.Error("error while retrieving database connection",
			log.WithField("name", name),
			log.WithError(ErrNotOpened),
		)
		// exit early
		panic(ErrNotOpened)
	}

	return instances[name]
}

func Open(s *settings.SqlDatabase, names ...string) error {
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
		return ErrAlreadyOpened
	}

	// TODO: add gorm configurations
	db, err := sql.Open(s.Driver, s.Uri)
	if err != nil {
		return err
	}

	db.SetMaxIdleConns(s.MaxIdleConns)
	db.SetMaxOpenConns(s.MaxOpenConns)
	db.SetConnMaxLifetime(s.ConnMaxLifetime)

	err = initORM(name, s, db)
	if err == ErrDialectNotSupported {
		log.Info("database doesn't support ORM", log.WithField("name", name))
	} else if err != nil && err != ErrDialectNotSupported {
		return err
	}

	instances[name] = db
	return nil
}

func Close() error {
	// close ORM instances first
	err := closeORM()
	if err != nil {
		return err
	}

	// keep try closes another db connection even when error
	// and return the last errors
	for name, db := range instances {
		// exclude ErrConnDone as it is possible already closed by ORM
		if e := db.Close(); e != nil && e != sql.ErrConnDone {
			log.Error("error when closing database",
				log.WithField("name", name),
				log.WithError(err))
			err = e // return the last error as errors
		}
	}

	return err
}
