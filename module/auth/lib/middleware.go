package lib

import (
	"github.com/euiko/webapp/core"
)

func AuthRequiredMiddleware(app core.App) core.MiddlewareFunc {
	return core.MustGetModule[Module](app).Middleware()
}
