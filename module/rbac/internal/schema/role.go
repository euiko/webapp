package schema

import (
	"github.com/euiko/webapp/db/sqldb"
	"github.com/uptrace/bun"
)

type (
	Role struct {
		bun.BaseModel `bun:"table:rbac.roles"`
		sqldb.BaseSchema

		ID          int64   `bun:"id,pk,autoincrement"`
		Name        string  `bun:"name,unique,notnull,nullzero"`
		PrettyName  string  `bun:"pretty_name,notnull,nullzero"`
		Description string  `bun:"description,nullzero"`
		Permissions []int64 `bun:"permissions,notnull,nullzero,array"`
	}
)
