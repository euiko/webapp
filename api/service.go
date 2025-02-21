package api

import "github.com/go-chi/chi/v5"

type (
	Service interface {
		Route(router chi.Router)
	}

	APIService interface {
		APIRoute(router chi.Router)
	}
)
