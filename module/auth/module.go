package auth

import (
	"context"
	"time"

	"github.com/euiko/webapp/api"
	"github.com/euiko/webapp/pkg/token"
	"github.com/euiko/webapp/settings"
)

type (
	Module struct {
		settings      Settings
		tokenEncoding token.Encoding
	}
)

func NewModule() *Module {
	return &Module{
		settings: Settings{
			Enabled: false,
			TokenEncoding: TokenEncodingSettings{
				Type:         "headless-jwt",
				JWTAlgorithm: "HS256",
				JWTIssuer:    "webapp",
				JWTAudience:  "webapp-server",
				JWTTimeout:   24 * time.Hour,
				HSKey:        "",
			},
		},
		tokenEncoding: nil,
	}
}

func (m *Module) DefaultSettings(s *settings.Settings) {
	s.SetExtra("auth", &m.settings)
}

func (m *Module) Init(ctx context.Context, s *settings.Settings) error {
	var err error

	m.tokenEncoding, err = NewTokenEncoding(s)
	if err != nil {
		return err
	}

	return nil
}

func (m *Module) Close() error {
	return nil
}

func ModuleFactory() func() api.Module {
	return func() api.Module {
		return NewModule()
	}
}
