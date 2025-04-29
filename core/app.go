package core

import (
	"context"

	"github.com/euiko/webapp/settings"
)

type (
	App interface {
		Settings() *settings.Settings
		Start(context.Context) error
		Modules() []Module
	}
)
