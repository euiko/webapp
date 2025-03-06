package cache

import (
	"context"
	"slices"
	"sync"
	"time"
)

type (
	InMemoryCache struct {
		mutex       sync.RWMutex
		data        map[string]cacheItem
		timeoutData map[int64][]string
		timeouts    []int64

		ctx    context.Context
		cancel func()
	}

	cacheItem struct {
		expireTime time.Time
		value      interface{}
	}
)

func NewInMemory() *InMemoryCache {
	cache := InMemoryCache{
		data:        make(map[string]cacheItem),
		timeoutData: make(map[int64][]string),
		mutex:       sync.RWMutex{},
		timeouts:    []int64{},
	}

	go cache.start()
	return &cache
}

func (c *InMemoryCache) Set(key string, value interface{}, opts ...SetOption) error {
	config := newSetConfig(opts...)

	item := cacheItem{
		value: value,
	}
	if config.Timeout > 0 {
		item.expireTime = time.Now().Add(config.Timeout)
	}

	return c.write(func() error {
		c.data[key] = item
		if !item.expireTime.IsZero() {
			c.deleteTimeout(key)
			c.setTimeout(key, item.expireTime)
		}

		return nil
	})
}

func (c *InMemoryCache) Get(key string) (interface{}, error) {
	var value interface{}
	err := c.read(func() error {
		item, ok := c.data[key]
		if !ok {
			return ErrKeyNotFound
		}

		value = item.value
		return nil
	})

	return value, err
}

func (c *InMemoryCache) Delete(key string) error {
	return c.write(func() error {
		// ensure key exists
		if _, err := c.getItem(key); err != nil {
			return err
		}

		c.deleteItem(key)
		return nil
	})
}

func (c *InMemoryCache) start() {
	for {
		ttl, ok := c.getFirstTimeout()
		c.startNewContext(ttl)

		// wait until context canceled
		<-c.ctx.Done()

		// directly skip delete keys when timeout still empty
		if !ok {
			continue
		}

		c.write(func() error {
			keys, ok := c.timeoutData[ttl.UnixMilli()]
			if !ok {
				return nil
			}

			for _, k := range keys {
				c.deleteItem(k)
			}

			c.clearTimeout(ttl)
			return nil
		})

	}
}

func (c *InMemoryCache) getItem(key string) (cacheItem, error) {
	item, ok := c.data[key]

	var err error
	if !ok {
		err = ErrKeyNotFound
	}

	return item, err
}

func (c *InMemoryCache) deleteItem(key string) {
	delete(c.data, key)
	c.deleteTimeout(key)
}

func (c *InMemoryCache) clearTimeout(t time.Time) {
	delete(c.timeoutData, t.UnixMilli())
	index, ok := slices.BinarySearch(c.timeouts, t.UnixMilli())
	if ok {
		c.timeouts = slices.Delete(c.timeouts, index, index+1)
	}
}

func (c *InMemoryCache) deleteTimeout(key string) {
	item, ok := c.data[key]
	if !ok || item.expireTime.IsZero() {
		// key not found or empty expireTime, do nothing
		return
	}

	expireMilli := item.expireTime.UnixMilli()
	keys := c.timeoutData[expireMilli]
	index := slices.Index(keys, key)

	// no index found, do nothing
	if index < 0 {
		return
	}

	c.timeoutData[expireMilli] = slices.Delete(keys, index, index+1)
	if len(c.timeoutData[expireMilli]) == 0 {
		delete(c.timeoutData, expireMilli)
	}

	arrIndex, ok := slices.BinarySearch(c.timeouts, expireMilli)
	if len(c.timeoutData[expireMilli]) == 0 && ok {
		c.timeouts = slices.Delete(c.timeouts, arrIndex, arrIndex+1)
		// once deleted cancel the current context
		if c.ctx != nil {
			c.cancel()
		}
	}
}

func (c *InMemoryCache) setTimeout(key string, expireAt time.Time) {

	keys := c.timeoutData[expireAt.UnixMilli()]
	keys = append(keys, key)
	c.timeoutData[expireAt.UnixMilli()] = keys

	// make sure it is ordered asc by using binary search
	arrIndex, _ := slices.BinarySearch(c.timeouts, expireAt.UnixMilli())
	newTimeouts := slices.Insert(c.timeouts, arrIndex, expireAt.UnixMilli())
	c.timeouts = newTimeouts

	// once done cancel the current context
	if c.ctx != nil {
		c.cancel()
	}
}

func (c *InMemoryCache) startNewContext(expireAt time.Time) {
	var (
		ctx    context.Context
		cancel func()
	)

	if expireAt.IsZero() {
		// use never expire context when expireAt is zero
		ctx, cancel = context.WithCancel(context.Background())
	} else {
		ctx, cancel = context.WithDeadline(context.Background(), expireAt)
	}

	c.write(func() error {
		c.ctx, c.cancel = ctx, cancel
		return nil
	})
}

func (c *InMemoryCache) getFirstTimeout() (time.Time, bool) {
	var (
		t  time.Time
		ok bool
	)

	c.read(func() error {
		if len(c.timeouts) == 0 {
			return nil
		}

		t = time.UnixMilli(c.timeouts[0])
		ok = true

		return nil
	})

	return t, ok
}

func (c *InMemoryCache) write(fn func() error) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return fn()
}

func (c *InMemoryCache) read(fn func() error) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return fn()
}
