package sqldb

import (
	"database/sql"
	"errors"
	"time"

	"github.com/euiko/webapp/pkg/log"
	"github.com/euiko/webapp/settings"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mssqldialect"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/schema"
)

type (
	BaseSchema struct {
		ID        int64     `bun:"id,pk,autoincrement"`
		CreatedAt time.Time `bun:"created_at,notnull,nullzero"`
		UpdatedAt time.Time `bun:"updated_at,notnull,nullzero"`
	}

	dialectFactory func(*settings.SqlDatabase) schema.Dialect

	OrmDB struct {
		*bun.DB
	}
)

var (
	ErrDialectAlreadyExists = errors.New("dialect already existed")
	ErrDialectNotSupported  = errors.New("no dialect supported for the selected driver")

	// for holds all the orm instances
	ormInstances    = make(map[string]*bun.DB)
	dialectRegistry = map[string]dialectFactory{
		"postgres": newPostgresDialect,
		"pgx":      newPostgresDialect,
		"mysql":    newMysqlDialect,
		"mssql":    newSqlServerDialect,
	}
)

func ORM(names ...string) OrmDB {
	name := defaultDbName
	if len(names) > 0 {
		name = names[0]
	}

	// ensure the database is opened
	if _, ok := ormInstances[name]; !ok {
		log.Error("error while retrieving ORM database instance",
			log.WithField("name", name),
			log.WithField("reason", "the database not yet opened or doesn't support ORM"),
		)
		// exit early
		panic(ErrNotOpened)
	}

	return OrmDB{ormInstances[name]}

}

func RegisterORMDialect(name string, f dialectFactory) error {
	if _, ok := dialectRegistry[name]; ok {
		return ErrDialectAlreadyExists
	}

	dialectRegistry[name] = f
	return nil
}

func initORM(name string, s *settings.SqlDatabase, db *sql.DB) error {
	dialectFactory, ok := dialectRegistry[s.Driver]
	if !ok {
		return ErrDialectNotSupported
	}

	ormInstances[name] = bun.NewDB(db, dialectFactory(s))
	return nil
}

func closeORM() error {
	var err error
	for name, db := range ormInstances {
		if e := db.Close(); e != nil {
			log.Error("error when closing orm database",
				log.WithField("name", name),
				log.WithError(err))
			err = e // return the last error as errors
		}
	}

	return err
}

func newPostgresDialect(s *settings.SqlDatabase) schema.Dialect {
	return pgdialect.New()
}

func newMysqlDialect(s *settings.SqlDatabase) schema.Dialect {
	return mysqldialect.New()
}

func newSqlServerDialect(s *settings.SqlDatabase) schema.Dialect {
	return mssqldialect.New()
}
