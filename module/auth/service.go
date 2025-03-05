package auth

import (
	"errors"
	"net/http"

	"github.com/euiko/webapp/api"
	"github.com/euiko/webapp/pkg/helper"
)

type (
	LoginPayload struct {
		LoginId  string `in:"form=login_id" json:"login_id" validate:"required"`
		Password string `in:"form=password" json:"password" validate:"required"`
	}

	LoginResponse struct {
		Token string `json:"token"`
	}
)

func (m *Module[User]) Route(r api.Router) {
	// only configure when it is enabled
	if m.settings.Enabled {
		m.initializeMiddleware(r)
	}
}

func (m *Module[User]) APIRoute(r api.Router) {
	// only configure it when it is enabled
	if !m.settings.Enabled {
		return
	}

	m.initializeMiddleware(r)

	// add the auth middleware
	api.PrivateRouter(r, func(r api.Router) {
		r.Post("/auth/logout", m.logoutHandler)
	})

	// public accessible routes
	r.Post("/auth/login", m.loginHandler)
}

func (m *Module[User]) initializeMiddleware(r api.Router) {
	api.PrivateRouter(r, func(r api.Router) {
		r.Use(NewMiddleware(m.tokenEncoding, nil))
	})
}

func (m *Module[User]) loginHandler(w http.ResponseWriter, r *http.Request) {
	var (
		payload LoginPayload
	)

	if err := helper.DecodeRequest(r, &payload); err != nil {
		helper.WriteResponse(w, err)
		return
	}

	// call before login hooks
	for _, hook := range m.hooks {
		if err := hook.BeforeLogin(r.Context(), payload.LoginId, payload.Password); err != nil {
			helper.WriteResponse(w, err)
			return
		}
	}

	user, err := m.userLoader.LoadUser(payload.LoginId, payload.Password)
	if err != nil {
		helper.WriteResponse(w, err)
		return
	}

	subject := (*user).Subject()
	token, err := m.tokenEncoding.Encode(subject, "webapp")
	if err != nil {
		helper.WriteResponse(w, err)
		return
	}

	// call after login hooks
	tokenStr := string(token)
	for _, hook := range m.hooks {
		if err := hook.AfterLogin(r.Context(), user, &tokenStr); err != nil {
			helper.WriteResponse(w, err)
			return
		}
	}

	response := LoginResponse{
		Token: tokenStr,
	}
	helper.WriteResponse(w, response)
}

func (m *Module[User]) logoutHandler(w http.ResponseWriter, r *http.Request) {
	if !IsAuthenticated(r.Context()) {
		helper.WriteResponse(w, errors.New("not authenticated"))
		return
	}

	// call before logout hooks
	for _, hook := range m.hooks {
		if err := hook.BeforeLogout(r.Context()); err != nil {
			helper.WriteResponse(w, err)
			return
		}
	}

	// TODO: revoke token

	// call after logout hooks
	for _, hook := range m.hooks {
		if err := hook.AfterLogout(r.Context()); err != nil {
			helper.WriteResponse(w, err)
			return
		}
	}

	helper.WriteResponse(
		w,
		map[string]interface{}{
			"message": "logout successful",
		},
	)
}
