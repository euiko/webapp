package auth

import "context"

type (
	Hook[User any] interface {
		BeforeLogin(ctx context.Context, loginId string, password string) error
		AfterLogin(ctx context.Context, user *User, token *string) error
		BeforeLogout(ctx context.Context) error
		AfterLogout(ctx context.Context) error
	}
)
