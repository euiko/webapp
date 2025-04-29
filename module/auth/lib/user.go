package lib

import (
	"context"

	"github.com/euiko/webapp/pkg/session"
)

type (
	User interface {
		session.Marshaller
		session.Unmarshaller

		LoginID() string
		Name() string
	}

	userContextKeyType struct{}
)

var (
	userContextKey = userContextKeyType{}
)

func CurrentUser(ctx context.Context) (User, bool) {
	v := ctx.Value(userContextKey)
	if v == nil {
		return nil, false
	}

	user, ok := v.(User)
	if !ok {
		return nil, false
	}

	return user, true
}

func CurrentUserTyped[T *User](ctx context.Context) (T, bool) {
	user, ok := CurrentUser(ctx)
	if !ok {
		return nil, false
	}

	typed, ok := user.(T)
	if !ok {
		return nil, false
	}

	return typed, true
}

func WithCurrentUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}
