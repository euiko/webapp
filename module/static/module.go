package static

import (
	"context"

	"github.com/euiko/webapp/api"
	"github.com/euiko/webapp/settings"
)

type (
	Module struct {
		settings *Settings
	}
)

func (m *Module) DefaultSettings(s *settings.Settings) {
	// set default settings
	s.Extra["static_server"] = Settings{
		Enabled: true,
		Embed: EmbedSettings{
			IndexPath: "index.html",
			UseMPA:    false,
		},
		Proxy: ProxySettings{
			Upstream: "http://localhost:5173",
		},
	}
}

func (m *Module) Init(ctx context.Context, s *settings.Settings) error {
	var err error
	m.settings, err = settings.GetExtra[Settings](s, "static_server")
	return err
}

func (m *Module) Close() error {
	return nil
}

func NewModule() api.Module {
	return &Module{}
}
