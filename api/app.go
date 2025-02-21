package api

import (
	"context"

	"github.com/euiko/webapp/settings"
)

type (
	App interface {
		Register(ModuleFactory)
		Run(context.Context) error
		Start(context.Context) error
	}

	SettingsLoader interface {
		DefaultSettings(s *settings.Settings)
	}
)
