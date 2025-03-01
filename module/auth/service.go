package auth

import (
	"net/http"

	"github.com/euiko/webapp/api"
	"github.com/euiko/webapp/pkg/log"
)

func (m *Module) Route(r api.Router) {
	// only configure when it is enabled
	if m.settings.Enabled {
		m.initializeMiddleware(r)
	}
}

func (m *Module) APIRoute(r api.Router) {
	// only configure it when it is enabled
	if !m.settings.Enabled {
		return
	}

	m.initializeMiddleware(r)

	// add the auth middleware
	api.PrivateRouter(r, func(r api.Router) {
		r.Get("/protected", m.protectedHandler)
	})

	// public accessible routes
	r.Route("/auth", func(r api.Router) {
		r.Get("/login", m.loginHandler)
	})
}

func (m *Module) initializeMiddleware(r api.Router) {
	api.PrivateRouter(r, func(r api.Router) {
		r.Use(NewMiddleware(m.tokenEncoding, nil))
	})
}

func (m *Module) loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("login handler")
	w.Write([]byte("login"))
}

func (m *Module) protectedHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("protected handler")
	w.Write([]byte("protected"))
}
