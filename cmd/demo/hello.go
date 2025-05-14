package main

import (
	"net/http"

	"github.com/euiko/webapp/core"
	authlib "github.com/euiko/webapp/module/auth/lib"
	"github.com/euiko/webapp/settings"
)

func newHelloService(app core.App) core.Module {
	return core.NewModule(
		core.ModuleWithAPIService(func(r core.Router, _ *settings.Settings) {
			r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("hello world!"))
			})

			authModule := core.MustGetModule[authlib.Module](app)
			r.Group(func(r core.Router) {
				r.Use(authModule.Middleware())

				r.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("hello protected world!"))
				})
			})
		}),
	)
}
