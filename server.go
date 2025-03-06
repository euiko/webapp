package webapp

import (
	"net/http"

	"github.com/euiko/webapp/api"
	"github.com/go-chi/chi/v5"
)

// internal createServer function
func (a *App) createServer() http.Server {
	// use chi as the router
	router := newRouter(chi.NewRouter())

	// use default middlewares
	router.Use(a.defaultMiddlewares...)

	// register routes
	router.Group(func(r api.Router) {
		for _, module := range a.modules {
			// register routes
			if service, ok := module.(api.Service); ok {
				service.Route(router)
			}
		}
	})

	// register api routes
	router.Route(a.settings.Server.ApiPrefix, func(r api.Router) {
		for _, module := range a.modules {
			// register routes
			if service, ok := module.(api.APIService); ok {
				service.APIRoute(r)
			}
		}
	})

	// creates http server
	// TODO: add https support
	return http.Server{
		Addr:         a.settings.Server.Addr,
		Handler:      router,
		ReadTimeout:  a.settings.Server.ReadTimeout,
		WriteTimeout: a.settings.Server.WriteTimeout,
		IdleTimeout:  a.settings.Server.IdleTimeout,
	}
}
