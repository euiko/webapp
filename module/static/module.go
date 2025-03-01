package static

import (
	"context"

	"github.com/euiko/webapp/api"
	"github.com/euiko/webapp/settings"
)

type (
	Module struct {
		settings Settings
	}
)

func (m *Module) DefaultSettings(s *settings.Settings) {
	// set default settings
	s.SetExtra("static_server", &m.settings)
}

func (m *Module) Init(ctx context.Context, s *settings.Settings) error {
	return nil
}

func (m *Module) Close() error {
	return nil
}

func NewModule() *Module {
	return &Module{
		settings: Settings{
			Enabled: true,
			Embed: EmbedSettings{
				IndexPath: "index.html",
				UseMPA:    false,
			},
			Proxy: ProxySettings{
				Upstream: "http://localhost:5173",
			},
		},
	}
}

func ModuleFactory() func() api.Module {
	return func() api.Module {
		return NewModule()
	}
}
