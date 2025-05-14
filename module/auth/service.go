package auth

import (
	"errors"
	"net/http"

	"github.com/euiko/webapp/core"
	"github.com/euiko/webapp/pkg/helper"
	"github.com/euiko/webapp/pkg/session"
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

func (m *Module[U]) APIRoute(r core.Router) {
	// only configure it when it is enabled
	if !m.settings.Enabled {
		return
	}

	r.With(m.Middleware()).Post("/auth/logout", m.logoutHandler)

	// public accessible routes
	r.Post("/auth/login", m.loginHandler)
}

func (m *Module[U]) loginHandler(w http.ResponseWriter, r *http.Request) {
	var (
		payload LoginPayload
		keys    = m.GetKeys()
	)

	if len(keys) == 0 {
		helper.WriteResponse(w, errors.New("invalid configuration"))
		return
	}

	if err := helper.DecodeRequestBody(r, &payload); err != nil {
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

	user, err := m.userLoader.LoadUser(r.Context(), payload.LoginId, payload.Password)
	if err != nil {
		helper.WriteResponse(w, err)
		return
	}

	subject := user.LoginID()
	key := keys[0] // use the first key to create token
	token, err := m.tokenEncoding.Encode(key, subject, "webapp")
	if err != nil {
		helper.WriteResponse(w, err)
		return
	}

	// write into session
	defer session.Add(r.Context(), "user", &user)

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

func (m *Module[U]) logoutHandler(w http.ResponseWriter, r *http.Request) {
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
