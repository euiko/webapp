package httpapi

import "github.com/euiko/webapp/pkg/common/base"

type (
	SearchParams struct {
		Keyword string `in:"query=keyword" json:"keyword"`
	}
)

func (s SearchParams) ToBase() base.SearchParams {
	return base.SearchParams{
		Keyword: s.Keyword,
	}
}

func (s *SearchParams) FromBase(base base.SearchParams) {
	s.Keyword = base.Keyword
}
