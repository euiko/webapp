package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/euiko/webapp/module/auth/lib"
	"github.com/euiko/webapp/pkg/helper"
	"github.com/euiko/webapp/pkg/log"
	"github.com/euiko/webapp/pkg/session"
	"github.com/euiko/webapp/pkg/token"
)

func newMiddleware[U lib.User](module *Module[U], unauthorizedHandler http.Handler) func(http.Handler) http.Handler {
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

			if len(authorization) == 0 {
				prohibited = true
			}

			// only try to decode the token if it is not prohibited
			if !prohibited {
				// try use all available keys
				keys := module.GetKeys()
				for _, key := range keys {
					token, err = module.TokenEncoding().Decode(key, []byte(authorization))
					if err == nil {
						break
					}
				}

				if err != nil {
					log.Error("failed to decode token", log.WithError(err))
					prohibited = true
				}
			}

			if prohibited {
				if unauthorizedHandler != nil {
					unauthorizedHandler.ServeHTTP(w, r)
				} else {
					helper.WriteResponse(w, errors.New("unauthorized"), helper.ResponseWithStatus(http.StatusUnauthorized))
				}
				return
			}

			var user lib.User
			err = session.Get(r.Context(), "user", &user)
			if err == session.ErrKeyNotFound {
				// load user from the user loader if not found in the session
				user, err = module.UserLoader().UserById(r.Context(), token.Subject)
				if err != nil {
					helper.WriteResponse(w, errors.New("internal server error"))
					return
				}

				session.Add(r.Context(), "user", user)
			} else if err != nil {
				// other errors
				helper.WriteResponse(w, errors.New("internal server error"))
				return
			}

			ctx := contextWithToken(r.Context(), token)
			ctx = lib.WithCurrentUser(ctx, user)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
