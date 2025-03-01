package main

import (
	"net/http"

	"github.com/euiko/webapp/api"
	"github.com/euiko/webapp/settings"
)

func newHelloService() api.Module {
	return api.NewModule(
		api.ModuleWithAPIService(func(r api.Router, _ *settings.Settings) {
			r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("hello world!"))
			})
		}),
	)
}
