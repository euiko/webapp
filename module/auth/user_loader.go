package auth

import (
	"context"

	"github.com/euiko/webapp/module/auth/lib"
)

type (
	userLoaderWrapped[U lib.User] struct {
		lib.UserLoader[U]
	}
)

func wrapUserLoader[U lib.User](u lib.UserLoader[U]) lib.UserLoader[lib.User] {
	return &userLoaderWrapped[U]{UserLoader: u}
}

func (l *userLoaderWrapped[U]) UserById(ctx context.Context, loginId string) (lib.User, error) {
	return l.UserLoader.UserById(ctx, loginId)
}

func (l *userLoaderWrapped[U]) LoadUser(ctx context.Context, loginId string, password string) (lib.User, error) {
	return l.UserLoader.LoadUser(ctx, loginId, password)
}
