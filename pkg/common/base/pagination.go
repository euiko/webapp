package base

type (
	PaginationParams struct {
		Offset int `validate:"gte=0"`
		Limit  int `validate:"gt=0,lte=1000"`
	}
)
