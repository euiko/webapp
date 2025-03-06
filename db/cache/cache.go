package cache

import (
	"errors"
	"sync"
	"time"
)

type (
	Cache interface {
		Set(key string, value interface{}, opts ...SetOption) error
		Get(key string) (interface{}, error)
		Delete(key string) error
	}

	SetConfig struct {
		Timeout time.Duration
	}

	SetOption func(c *SetConfig)
)

var (
	inMemoryCache     Cache
	inMemoryCacheOnce sync.Once

	ErrKeyNotFound = errors.New("key not found")
)

func SetWithTimeout(timeout time.Duration) SetOption {
	return func(c *SetConfig) {
		c.Timeout = timeout
	}
}

func InMemory() Cache {
	inMemoryCacheOnce.Do(func() {
		inMemoryCache = NewInMemory()
	})

	return inMemoryCache
}

func newSetConfig(opts ...SetOption) *SetConfig {
	c := SetConfig{
		Timeout: 0,
	}

	for _, opt := range opts {
		opt(&c)
	}

	return &c
}
