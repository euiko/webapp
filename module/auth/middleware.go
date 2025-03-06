package auth

import (
	"net/http"
	"strings"

	"github.com/euiko/webapp/pkg/log"
	"github.com/euiko/webapp/pkg/token"
)

func (m *Module[User]) newMiddleware(unauthorizedHandler http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				prohibited bool
				token      *token.Token
				err        error
			)

			authorization := r.Header.Get("Authorization")
			if !strings.HasPrefix(authorization, "Bearer") {
				prohibited = true
			}

			// clean up the authorization header to obtain token
			authorization = strings.TrimPrefix(authorization, "Bearer")
			authorization = strings.TrimSpace(authorization)

			// try use all available keys
			keys := m.getKeys()
			for _, key := range keys {
				token, err = m.tokenEncoding.Decode(key, []byte(authorization))
				if err == nil {
					break
				}
			}

			if err != nil {
				log.Error("failed to decode token", log.WithError(err))
				prohibited = true
			}

			if prohibited {
				if unauthorizedHandler != nil {
					unauthorizedHandler.ServeHTTP(w, r)
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("unauthorized"))
				}
				return
			}

			ctx := contextWithToken(r.Context(), token)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
