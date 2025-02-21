package api

import (
	"context"

	"github.com/euiko/webapp/settings"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
)

type (
	Module interface {
		Init(context.Context, *settings.Settings) error
		Close() error
	}

	ModuleFactory func() Module

	ModuleOption func(*module)

	module struct {
		settings           *settings.Settings
		initFunc           func(ctx context.Context, s *settings.Settings) error
		closeFunc          func(s *settings.Settings) error
		settingsLoaderFunc func(s *settings.Settings)
		serviceFunc        func(router chi.Router, s *settings.Settings)
		apiServiceFunc     func(router chi.Router, s *settings.Settings)
		cliFunc            func(cmd *cobra.Command, s *settings.Settings)
	}
)

func ModuleWithInit(f func(context.Context, *settings.Settings) error) ModuleOption {
	return func(m *module) {
		m.initFunc = f
	}
}

func ModuleWithClose(f func(*settings.Settings) error) ModuleOption {
	return func(m *module) {
		m.closeFunc = f
	}
}

func ModuleWithSettingsLoader(f func(*settings.Settings)) ModuleOption {
	return func(m *module) {
		m.settingsLoaderFunc = f
	}
}

func ModuleWithService(f func(chi.Router, *settings.Settings)) ModuleOption {
	return func(m *module) {
		m.serviceFunc = f
	}
}

func ModuleWithAPIService(f func(chi.Router, *settings.Settings)) ModuleOption {
	return func(m *module) {
		m.apiServiceFunc = f
	}
}

func ModuleWithCLI(f func(*cobra.Command, *settings.Settings)) ModuleOption {
	return func(m *module) {
		m.cliFunc = f
	}
}

func NewModule(opts ...ModuleOption) Module {
	m := &module{}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m *module) Init(ctx context.Context, s *settings.Settings) error {
	m.settings = s
	if m.initFunc != nil {
		return m.initFunc(ctx, s)
	}
	return nil
}

func (m *module) Close() error {
	if m.closeFunc != nil {
		return m.closeFunc(m.settings)
	}
	return nil
}

func (m *module) DefaultSettings(s *settings.Settings) {
	if m.settingsLoaderFunc != nil {
		m.settingsLoaderFunc(s)
	}
}

func (m *module) Route(router chi.Router) {
	if m.serviceFunc != nil {
		m.serviceFunc(router, m.settings)
	}
}

func (m *module) APIRoute(router chi.Router) {
	if m.apiServiceFunc != nil {
		m.apiServiceFunc(router, m.settings)
	}
}

func (m *module) Command(cmd *cobra.Command) {
	if m.cliFunc != nil {
		m.cliFunc(cmd, m.settings)
	}
}
