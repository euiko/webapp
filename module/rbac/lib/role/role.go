package role

import "time"

type (
	Base struct {
		Name        string `validate:"required"`
		PrettyName  string
		Description string
		Permissions []int64 `validate:"required"`
	}

	New Base

	Update struct {
		PrettyName  string
		Description string
		Permissions []int64 `validate:"required"`
	}

	Role struct {
		Base
		ID        int64
		CreatedAt time.Time
		UpdatedAt time.Time
	}
)
