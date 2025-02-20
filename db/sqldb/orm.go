package sqldb

import (
	"database/sql"
	"errors"

	"github.com/euiko/webapp/pkg/log"
	"github.com/euiko/webapp/settings"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type (
	dialectFactory func(*settings.SqlDatabase, *sql.DB) gorm.Dialector
)

var (
	ErrDialectAlreadyExists = errors.New("dialect already existed")
	ErrDialectNotSupported  = errors.New("no dialect supported for the selected driver")

	// for holds all the orm instances
	ormInstances    = make(map[string]*gorm.DB)
	dialectRegistry = map[string]dialectFactory{
		"postgres": newPostgresDialect,
		"pgx":      newPostgresDialect,
		"mysql":    newMysqlDialect,
		"azuersql": newSqlServerDialect,
	}
)

func ORM(names ...string) *gorm.DB {
	name := defaultDbName
	if len(names) > 0 {
		name = names[0]
	}

	// ensure the database is opened
	if _, ok := ormInstances[name]; !ok {
		log.Error(errNotOpened.Error(),
			log.WithField("name", name),
			log.WithError(errNotOpened),
		)
		// exit early
		panic(errNotOpened)
	}

	return ormInstances[name]

}

func RegisterORMDialect(name string, f dialectFactory) error {
	if _, ok := dialectRegistry[name]; ok {
		return ErrDialectAlreadyExists
	}

	dialectRegistry[name] = f
	return nil
}

func initORM(name string, s *settings.SqlDatabase, db *sql.DB) error {
	var config gorm.Config

	dialectFactory, ok := dialectRegistry[s.Driver]
	if !ok {
		return ErrDialectNotSupported
	}

	gormDb, err := gorm.Open(dialectFactory(s, db), &config)
	if err != nil {
		return err
	}

	ormInstances[name] = gormDb
	return nil
}

func newPostgresDialect(s *settings.SqlDatabase, db *sql.DB) gorm.Dialector {
	return postgres.New(postgres.Config{Conn: db})
}

func newMysqlDialect(s *settings.SqlDatabase, db *sql.DB) gorm.Dialector {
	return mysql.New(mysql.Config{Conn: db})
}

func newSqlServerDialect(s *settings.SqlDatabase, db *sql.DB) gorm.Dialector {
	return sqlserver.New(sqlserver.Config{Conn: db})
}
