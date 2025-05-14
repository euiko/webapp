package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/euiko/webapp/core"
	"github.com/euiko/webapp/db/cache"
	"github.com/euiko/webapp/pkg/helper"
	"github.com/euiko/webapp/pkg/log"
	"github.com/euiko/webapp/pkg/token"
	"github.com/euiko/webapp/settings"

	"github.com/euiko/webapp/module/auth/lib"
)

type (
	Module[U lib.User] struct {
		app                 core.App
		settings            Settings
		tokenEncoding       token.Encoding
		userLoader          lib.UserLoader[U]
		hooks               []lib.Hook[U]
		keyStore            token.KeyStore
		middleware          func(http.Handler) http.Handler
		unauthorizedHandler http.Handler
	}

	ModuleOption[U lib.User] func(*Module[U])
)

const (
	cacheKeyKeys = "auth:keys"
)

func WithUnauthorizedHandler[U lib.User](handler http.Handler) ModuleOption[U] {
	return func(m *Module[U]) {
		m.unauthorizedHandler = handler
	}
}

func ModuleFactory[U lib.User](
	userLoader lib.UserLoader[U],
	options ...ModuleOption[U],
) core.ModuleFactory {
	return func(app core.App) core.Module {
		return NewModule[U](app, userLoader, options...)
	}
}

func NewModule[U lib.User](
	app core.App,
	userLoader lib.UserLoader[U],
	options ...ModuleOption[U],
) *Module[U] {
	m := Module[U]{
		app: app,
		settings: Settings{
			Enabled: false,
			TokenEncoding: TokenEncodingSettings{
				Type:         "headless-jwt",
				JWTAlgorithm: "HS256",
				JWTIssuer:    "webapp",
				JWTAudience:  "webapp",
				JWTTimeout:   24 * time.Hour,
				Keys: []string{
					helper.EncodeBase64(helper.Hash([]byte("secret"), helper.HashSHA256)),
				},
			},
		},
		tokenEncoding: nil,
		userLoader:    userLoader,
	}

	for _, opt := range options {
		opt(&m)
	}

	return &m
}

func (m *Module[U]) DefaultSettings(s *settings.Settings) {
	s.SetExtra("auth", &m.settings)
}

func (m *Module[U]) Init(ctx context.Context, s *settings.Settings) error {
	if !m.settings.Enabled {
		return nil
	}

	var err error
	m.tokenEncoding, err = NewTokenEncoding(&m.settings)
	if err != nil {
		return err
	}

	// add keys in configurations
	for _, key := range m.settings.TokenEncoding.Keys {
		m.keyStore.Add(key, token.NewSymetricKey([]byte(key)))
	}

	return nil
}

func (m *Module[U]) Close() error {
	return nil
}

func (m *Module[U]) GetKeys() []token.Key {
	cached, err := cache.InMemory().Get(cacheKeyKeys)
	if err == cache.ErrKeyNotFound {
		keys := m.keyStore.Keys()
		cache.InMemory().Set(cacheKeyKeys, keys, cache.SetWithTimeout(10*time.Minute))
		return keys
	}

	if err != nil {
		log.Error("failed when retrieving keys from cache", log.WithError(err))
		return nil
	}

	return cached.([]token.Key)
}

func (m *Module[U]) TokenEncoding() token.Encoding {
	return m.tokenEncoding
}

func (m *Module[U]) UserLoader() lib.UserLoader[lib.User] {
	return wrapUserLoader(m.userLoader)
}

func (m *Module[U]) Middleware() core.MiddlewareFunc {
	if m.middleware == nil {
		m.middleware = newMiddleware(m, nil)
	}

	return m.middleware
}
