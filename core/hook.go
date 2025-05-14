package core

import (
	"context"

	"github.com/euiko/webapp/settings"
)

type (
	ModuleLoadedHook interface {
		ModuleLoaded(Module)
	}

	PostRouterHook interface {
		PostRoute(Router)
	}

	SettingsLoaderHook interface {
		DefaultSettings(*settings.Settings)
	}

	BeforeStartHook interface {
		BeforeStart(context.Context) error
	}
)
