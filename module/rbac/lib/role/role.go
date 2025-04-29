package role

import "time"

type (
	BaseRole struct {
		Name        string
		PrettyName  string
		Description string
		Permissions []int64
	}

	NewRole BaseRole

	Role struct {
		BaseRole
		ID        int64
		CreatedAt time.Time
		UpdatedAt time.Time
	}
)
