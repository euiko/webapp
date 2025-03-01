package auth

import (
	"net/http"
	"strings"

	"github.com/euiko/webapp/pkg/log"
	"github.com/euiko/webapp/pkg/token"
)

func NewMiddleware(tokenEncoding token.Encoding, unauthorizedHandler http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				prohibited bool
			)

			authorization := r.Header.Get("Authorization")
			if !strings.HasPrefix(authorization, "Bearer") {
				prohibited = true
			}

			// clean up the authorization header to obtain token
			authorization = strings.TrimPrefix(authorization, "Bearer")
			authorization = strings.TrimSpace(authorization)

			token, err := tokenEncoding.Decode([]byte(authorization))
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
