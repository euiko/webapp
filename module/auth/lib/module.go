package lib

import (
	"github.com/euiko/webapp/core"
	"github.com/euiko/webapp/pkg/token"
)

type (
	Module interface {
		GetKeys() []token.Key
		UserLoader() UserLoader[User]
		TokenEncoding() token.Encoding
		Middleware() core.MiddlewareFunc
	}
)
