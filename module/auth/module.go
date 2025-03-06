package auth

import (
	"context"
	"time"

	"github.com/euiko/webapp/api"
	"github.com/euiko/webapp/db/cache"
	"github.com/euiko/webapp/pkg/helper"
	"github.com/euiko/webapp/pkg/log"
	"github.com/euiko/webapp/pkg/token"
	"github.com/euiko/webapp/settings"
)

type (
	Module[User Tokenable] struct {
		settings      Settings
		tokenEncoding token.Encoding
		userLoader    UserLoader[User]
		hooks         []Hook[User]
		keyStore      token.KeyStore
	}

	Tokenable interface {
		Subject() string
	}
)

const (
	cacheKeyKeys = "auth:keys"
)

func ModuleFactory[User Tokenable](
	userLoader UserLoader[User],
	hooks ...Hook[User],
) func() api.Module {
	return func() api.Module {
		return NewModule[User](userLoader)
	}
}

func NewModule[User Tokenable](
	userLoader UserLoader[User],
	hooks ...Hook[User],
) *Module[User] {
	return &Module[User]{
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
}

func (m *Module[User]) DefaultSettings(s *settings.Settings) {
	s.SetExtra("auth", &m.settings)
}

func (m *Module[User]) Init(ctx context.Context, s *settings.Settings) error {
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

func (m *Module[User]) Close() error {
	return nil
}

func (m *Module[User]) getKeys() []token.Key {
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
