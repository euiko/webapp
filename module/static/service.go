package static

import "github.com/go-chi/chi/v5"

func (m *Module) Route(router chi.Router) {
	if m.settings.Enabled {
		createStaticRoutes(router, m.settings)
	}
}
