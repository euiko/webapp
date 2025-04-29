package static

import "github.com/euiko/webapp/core"

func (m *Module) Route(router core.Router) {
	if m.settings.Enabled {
		// create another router without auth middleware
		createStaticRoutes(router, &m.settings)
	}
}
