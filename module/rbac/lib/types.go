package lib

import "github.com/euiko/webapp/pkg/common/base"

type (
	ListAllRolesParams struct {
		base.SearchParams
		base.PaginationParams
	}
)
