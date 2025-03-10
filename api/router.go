package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type (
	// Router is an interface that provide routing methods that mostly
	// resemble from chi.Router with difference in:
	// - Additional NamedGroup function
	Router interface {
		http.Handler
		chi.Routes

		// Use appends one or more middlewares onto the Router stack.
		Use(middlewares ...func(http.Handler) http.Handler)

		// With adds inline middlewares for an endpoint handler.
		With(middlewares ...func(http.Handler) http.Handler) Router

		// Group adds a new inline-Router along the current routing
		// path, with a fresh middleware stack for the inline-Router.
		Group(fn func(r Router)) Router

		// Route mounts a sub-Router along a `pattern`` string.
		Route(pattern string, fn func(r Router)) Router

		// Mount attaches another http.Handler along ./pattern/*
		Mount(pattern string, h http.Handler)

		// Handle and HandleFunc adds routes for `pattern` that matches
		// all HTTP methods.
		Handle(pattern string, h http.Handler)
		HandleFunc(pattern string, h http.HandlerFunc)

		// Method and MethodFunc adds routes for `pattern` that matches
		// the `method` HTTP method.
		Method(method, pattern string, h http.Handler)
		MethodFunc(method, pattern string, h http.HandlerFunc)

		// HTTP-method routing along `pattern`
		Connect(pattern string, h http.HandlerFunc)
		Delete(pattern string, h http.HandlerFunc)
		Get(pattern string, h http.HandlerFunc)
		Head(pattern string, h http.HandlerFunc)
		Options(pattern string, h http.HandlerFunc)
		Patch(pattern string, h http.HandlerFunc)
		Post(pattern string, h http.HandlerFunc)
		Put(pattern string, h http.HandlerFunc)
		Trace(pattern string, h http.HandlerFunc)

		// NotFound defines a handler to respond whenever a route could
		// not be found.
		NotFound(h http.HandlerFunc)

		// MethodNotAllowed defines a handler to respond whenever a method is
		// not allowed.
		MethodNotAllowed(h http.HandlerFunc)

		// NamedGroup adds or reuse inline-Router that unique by its name
		// along the current routing path, with a fresh middleware stack for the inline-Router.
		NamedGroup(name string, fn func(r Router)) Router
	}
)

const (
	protectedGroupName = "private"
)

// PrivateRouter returns a group/sub-router for private/protected routes
// this can be used to add protected resources to the router
func PrivateRouter(r Router, fn func(r Router)) Router {
	return r.NamedGroup(protectedGroupName, fn)
}
