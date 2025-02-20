package api

import (
	"context"

	"github.com/euiko/webapp/settings"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
)

type (
	Module interface {
		Init(ctx context.Context) error
		Close() error
	}

	ModuleFactory func(*settings.Settings) Module

	APIService interface {
		APIRoute(router chi.Router)
	}

	CLI interface {
		Command(cmd *cobra.Command)
	}

	ModuleOption func(*module)

	module struct {
		initFunc       func(ctx context.Context) error
		closeFunc      func() error
		apiServiceFunc func(router chi.Router)
		cliFunc        func(cmd *cobra.Command)
	}
)

func ModuleWithInit(f func(ctx context.Context) error) ModuleOption {
	return func(m *module) {
		m.initFunc = f
	}
}

func ModuleWithClose(f func() error) ModuleOption {
	return func(m *module) {
		m.closeFunc = f
	}
}

func ModuleWithAPIService(f func(router chi.Router)) ModuleOption {
	return func(m *module) {
		m.apiServiceFunc = f
	}
}

func ModuleWithCLI(f func(cmd *cobra.Command)) ModuleOption {
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

func (m *module) Init(ctx context.Context) error {
	if m.initFunc != nil {
		return m.initFunc(ctx)
	}
	return nil
}

func (m *module) Close() error {
	if m.closeFunc != nil {
		return m.closeFunc()
	}
	return nil
}

func (m *module) APIRoute(router chi.Router) {
	if m.apiServiceFunc != nil {
		m.apiServiceFunc(router)
	}
}
func (m *module) Command(cmd *cobra.Command) {
	if m.cliFunc != nil {
		m.cliFunc(cmd)
	}
}
