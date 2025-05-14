package lib

import "context"

type (
	UserLoaderFactory[U User] func(context.Context) UserLoader[U]
	UserLoader[U User]        interface {
		UserById(ctx context.Context, loginId string) (U, error)
		LoadUser(ctx context.Context, loginId string, password string) (U, error)
	}
)
