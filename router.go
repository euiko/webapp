package webapp

import (
	"net/http"

	"github.com/euiko/webapp/api"
	"github.com/go-chi/chi/v5"
)

type (
	router struct {
		chi.Router

		namedGroups map[string]api.Router
	}
)

// With adds inline middlewares for an endpoint handler.
func (r *router) With(middlewares ...func(http.Handler) http.Handler) api.Router {
	return newRouter(r.Router.With(middlewares...))
}

// Group adds a new inline-Router along the current routing
// path, with a fresh middleware stack for the inline-Router.
func (r *router) Group(fn func(r api.Router)) api.Router {
	group := r.Router.Group(func(r chi.Router) {
		if fn != nil {
			fn(newRouter(r))
		}
	})

	return newRouter(group)
}

// Route mounts a sub-Router along a `pattern“ string.
func (r *router) Route(pattern string, fn func(r api.Router)) api.Router {
	subRouter := r.Router.Route(pattern, func(r chi.Router) {
		if fn != nil {
			fn(newRouter(r))
		}
	})

	return newRouter(subRouter)
}

// NamedGroup adds or reuse inline-Router that unique by its name
// along the current routing path, with a fresh middleware stack for the inline-Router.
func (r *router) NamedGroup(name string, fn func(r api.Router)) api.Router {
	router, ok := r.namedGroups[name]
	if !ok {
		router = r.Group(nil)
		r.namedGroups[name] = router
	}

	if fn != nil {
		fn(router)
	}
	return router
}

func newRouter(r chi.Router) api.Router {
	return &router{
		Router:      r,
		namedGroups: make(map[string]api.Router),
	}
}
