package webapp

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/euiko/webapp/api"
	"github.com/euiko/webapp/db"
	"github.com/euiko/webapp/internal/cli"
	"github.com/euiko/webapp/pkg/log"
	"github.com/euiko/webapp/pkg/signal"
	"github.com/euiko/webapp/settings"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
)

type (
	App struct {
		name      string
		shortName string
		settings  settings.Settings

		registry           []api.ModuleFactory
		modules            []api.Module
		defaultMiddlewares []Middleware
	}

	Middleware func(http.Handler) http.Handler

	Option func(*App)
)

func WithDefaultMiddlewares(middlewares ...Middleware) Option {
	return func(a *App) {
		a.defaultMiddlewares = middlewares
	}
}

func New(name string, shortName string, opts ...Option) *App {
	app := App{
		name:               name,
		shortName:          shortName,
		modules:            []api.Module{},
		defaultMiddlewares: []Middleware{},
		settings:           settings.New(),
	}

	// apply options
	for _, opt := range opts {
		opt(&app)
	}

	// add built-in modules
	app.registry = append(app.registry, app.builtInModules()...)

	return &app
}

// Register a module factory function to the app
func (a *App) Register(f api.ModuleFactory) {
	a.registry = append(a.registry, f)
}

// Run the app
func (a *App) Run(ctx context.Context) error {
	// instantiate modules
	log.Trace("instantiating modules...")
	a.modules = make([]api.Module, len(a.registry))
	for i, factory := range a.registry {
		a.modules[i] = factory()
	}

	// load settings
	// TODO: handle error when unmarshaling
	unmarshaler, err := a.loadSettings()
	if err != nil && err != ErrConfigNotFound {
		return err
	}

	if unmarshaler != nil {
		if err := unmarshaler.Unmarshal(&a.settings); err != nil {
			return err
		}
	}

	// initialize logger
	initializeLogger(a.settings.Log)

	// initialize modules
	log.Trace("initializing modules...")
	for _, module := range a.modules {
		module.Init(ctx, &a.settings)
	}

	rootCmd := a.initializeCli()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		return err
	}

	// close all modules
	log.Trace("closing modules...")
	for _, module := range a.modules {
		module.Close()
	}
	return nil
}

func (a *App) Start(ctx context.Context) error {
	var err error
	// create and initialize server
	log.Info("starting the server...", log.WithField("addr", a.settings.Server.Addr))
	server := a.createServer()
	if err := db.Init(&a.settings.DB); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		if e := server.ListenAndServe(); e != nil && e != http.ErrServerClosed {
			err = e
			cancel()
		}
	}()

	// wait for signal to be done
	signal := signal.NewSignalNotifier()
	signal.OnSignal(func(ctx context.Context, sig os.Signal) bool {
		return true // exit on receiving any signal
	})
	signal.Wait(ctx)

	// when the err is being set means there is error on ListenAndServe
	if err != nil {
		return err
	}

	// close the server within 120s
	log.Info("closing the server...")
	defer db.Close()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer shutdownCancel() // ensure no context leak on graceful shutdown
	return server.Shutdown(shutdownCtx)
}

func (a *App) initializeCli() *cobra.Command {
	rootCmd := cobra.Command{
		Use: a.name,
	}

	for m := range a.modules {
		if cli, ok := a.modules[m].(api.CLI); ok {
			cli.Command(&rootCmd)
		}
	}

	return &rootCmd
}

// internal createServer function
func (a *App) createServer() http.Server {
	// use chi as the router
	router := chi.NewRouter()

	// use default middlewares
	for _, middleware := range a.defaultMiddlewares {
		router.Use(middleware)
	}

	// register routes
	for _, module := range a.modules {
		// register routes
		if service, ok := module.(api.Service); ok {
			service.Route(router)
		}
	}

	// register api routes
	router.Route(a.settings.Server.ApiPrefix, func(r chi.Router) {
		for _, module := range a.modules {
			// register routes
			if service, ok := module.(api.APIService); ok {
				service.APIRoute(r)
			}
		}
	})

	// creates http server
	// TODO: add https support
	return http.Server{
		Addr:         a.settings.Server.Addr,
		Handler:      router,
		ReadTimeout:  a.settings.Server.ReadTimeout,
		WriteTimeout: a.settings.Server.WriteTimeout,
		IdleTimeout:  a.settings.Server.IdleTimeout,
	}
}

func (a *App) builtInModules() []api.ModuleFactory {
	return []api.ModuleFactory{
		cli.Server(a),
		cli.Migration,
		cli.Settings,
	}
}

func initializeLogger(settings settings.Log) {
	// use LogrusLogger as default logger
	level := log.ParseLevel(settings.Level)
	log.SetDefault(log.NewLogrusLogger(level))
}
