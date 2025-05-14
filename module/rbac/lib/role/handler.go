package role

import "net/http"

type (
	handler struct {
		http.Handler
		permission *Permission
	}
)

func (h *handler) Permission() *Permission {
	return h.permission
}

func Handler(permission *Permission, next http.Handler) http.Handler {
	return &handler{
		Handler:    next,
		permission: permission,
	}
}
