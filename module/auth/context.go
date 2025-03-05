package auth

import (
	"context"

	"github.com/euiko/webapp/pkg/token"
)

type (
	tokenContextKeyType struct{}
	userContextKeyType  struct{}
)

var (
	tokenContextKey = tokenContextKeyType{}
	userContextKey  = userContextKeyType{}
)

func IsAuthenticated(ctx context.Context) bool {
	_, ok := TokenFromContext(ctx)
	return ok
}

func TokenFromContext(ctx context.Context) (*token.Token, bool) {
	v := ctx.Value(tokenContextKey)
	if v == nil {
		return nil, false
	}

	token, ok := v.(*token.Token)
	if !ok {
		return nil, false
	}

	return token, true
}

func UserFromContext[User any](ctx context.Context) (*User, bool) {
	v := ctx.Value(userContextKey)
	if v == nil {
		return nil, false
	}

	user, ok := v.(*User)
	if !ok {
		return nil, false
	}

	return user, true
}

func contextWithToken(ctx context.Context, token *token.Token) context.Context {
	return context.WithValue(ctx, tokenContextKey, token)
}

func contextWithUser[User any](ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}
