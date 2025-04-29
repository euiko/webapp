package role

import (
	"github.com/euiko/webapp/pkg/helper"
)

func ID(group, name string) int64 {
	b := []byte(group)
	b = append(b, byte(0))
	b = append(b, []byte(name)...)
	// use hash32 internally so it can be stored in int64 safely
	return int64(helper.Hash32(b))
}
