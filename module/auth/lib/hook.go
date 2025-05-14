package lib

import "context"

type (
	Hook[U User] interface {
		BeforeLogin(ctx context.Context, loginId string, password string) error
		AfterLogin(ctx context.Context, user U, token *string) error
		BeforeLogout(ctx context.Context) error
		AfterLogout(ctx context.Context) error
	}
)
