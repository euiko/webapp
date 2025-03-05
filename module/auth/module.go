package auth

import (
	"context"
	"time"

	"github.com/euiko/webapp/api"
	"github.com/euiko/webapp/pkg/token"
	"github.com/euiko/webapp/settings"
)

type (
	Module[User Tokenable] struct {
		settings      Settings
		tokenEncoding token.Encoding
		userLoader    UserLoader[User]
		hooks         []Hook[User]
	}

	Tokenable interface {
		Subject() string
	}
)

func NewModule[User Tokenable](
	userLoader UserLoader[User],
	hooks ...Hook[User],
) *Module[User] {
	return &Module[User]{
		settings: Settings{
			Enabled: false,
			TokenEncoding: TokenEncodingSettings{
				Type:         "headless-jwt",
				JWTAlgorithm: "HS256",
				JWTIssuer:    "webapp",
				JWTAudience:  "webapp",
				JWTTimeout:   24 * time.Hour,
				HSKey:        "",
			},
		},
		tokenEncoding: nil,
		userLoader:    userLoader,
	}
}

func (m *Module[User]) DefaultSettings(s *settings.Settings) {
	s.SetExtra("auth", &m.settings)
}

func (m *Module[User]) Init(ctx context.Context, s *settings.Settings) error {
	var err error

	m.tokenEncoding, err = NewTokenEncoding(&m.settings)
	if err != nil {
		return err
	}

	return nil
}

func (m *Module[User]) Close() error {
	return nil
}

func ModuleFactory[User Tokenable](
	userLoader UserLoader[User],
	hooks ...Hook[User],
) func() api.Module {
	return func() api.Module {
		return NewModule[User](userLoader)
	}
}
