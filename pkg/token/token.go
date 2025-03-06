package token

import (
	"time"
)

type (
	Token struct {
		Issuer    string
		Subject   string
		Audience  []string
		ExpiresAt time.Time
		IssuedAt  time.Time
	}
)
