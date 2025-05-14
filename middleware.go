package webapp

import (
	"net/http"

	"github.com/euiko/webapp/core"
	"github.com/euiko/webapp/pkg/log"
	"github.com/euiko/webapp/pkg/session"
	"github.com/euiko/webapp/settings"
)

func newInjectAppMiddleware(app core.App) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := contextWithApp(r.Context(), app)
			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
	}
}

func newSessionMiddleware(_ *settings.Settings) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			e := session.NewHTTPCookieEncoding(w, r)
			sessionValue, err := e.Decode()
			if err != nil {
				// just log the error
				log.Debug("failed to decode session, falling back to empty session", log.WithError(err))
				sessionValue = session.New()
			}

			ctx := session.WithContext(r.Context(), sessionValue)
			r = r.WithContext(ctx)

			h.ServeHTTP(w, r)

			e.Encode(sessionValue)
		})
	}
}
