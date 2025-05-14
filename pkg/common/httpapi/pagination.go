package httpapi

import "github.com/euiko/webapp/pkg/common/base"

type (
	PaginationParams struct {
		Page     int `in:"query=page" json:"page" validate:"gt=0"`
		PageSize int `in:"query=page_size" json:"page_size" validate:"gt=0,lte=1000"`
	}

	Pagination struct {
		Page     int `json:"page"`
		PageSize int `json:"page_size"`
		Total    int `json:"total"`
	}
)

func (p PaginationParams) ToBase() base.PaginationParams {
	return base.PaginationParams{
		Offset: (p.Page - 1) * p.PageSize,
		Limit:  p.PageSize,
	}
}
