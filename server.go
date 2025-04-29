package webapp

import (
	"net/http"

	"github.com/euiko/webapp/core"
	"github.com/go-chi/chi/v5"
)

// internal createServer function
func (a *App) createServer() http.Server {
	// use chi as the router
	router := newRouter(chi.NewRouter())

	// use default middlewares
	router.Use(newInjectAppMiddleware(a))
	router.Use(a.defaultMiddlewares...)
	router.Use(newSessionMiddleware(&a.settings))

	// register routes
	visitModules(a.modules, func(module core.ServiceModule) error {
		module.Route(router)
		return nil
	})
	// register api routes
	router.Route(a.settings.Server.ApiPrefix, func(r core.Router) {
		_ = visitModules(a.modules, func(module core.APIServiceModule) error {
			module.APIRoute(r)
			return nil
		})
	})

	// call post route hook
	_ = visitModules(a.modules, func(module core.PostRouterHook) error {
		module.PostRoute(router)
		return nil
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
