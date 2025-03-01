package static

import (
	"github.com/euiko/webapp/api"
)

func (m *Module) Route(router api.Router) {
	if m.settings.Enabled {
		// create another router without auth middleware
		createStaticRoutes(router, &m.settings)
	}
}
