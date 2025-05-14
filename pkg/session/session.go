package session

import (
	"context"
	"errors"
	"sync"

	"github.com/go-viper/mapstructure/v2"
)

type (
	Session struct {
		sync.Map
	}

	Marshaller interface {
		MarshalSession() (interface{}, error)
	}

	Unmarshaller interface {
		UnmarshalSession(interface{}) error
	}

	contextKey struct{}
)

var (
	sessionContextKey contextKey
	ErrNotInitialized = errors.New("session not initialized")
	ErrKeyNotFound    = errors.New("session not initialized")
)

func New() *Session {
	return new(Session)
}

func Add(ctx context.Context, key string, value interface{}) error {
	session, ok := fromContext(ctx)
	if !ok {
		return ErrNotInitialized
	}

	var err error
	if m, ok := value.(Marshaller); ok {
		value, err = m.MarshalSession()
	}

	if err != nil {
		return err
	}

	session.Store(key, value)
	return nil
}

func Delete(ctx context.Context, key string) error {
	session, ok := fromContext(ctx)
	if !ok {
		return ErrNotInitialized
	}

	session.Delete(key)
	return nil
}

func Get(ctx context.Context, key string, output interface{}) error {
	session, ok := fromContext(ctx)
	if !ok {
		return ErrNotInitialized
	}

	value, ok := session.Load(key)
	if !ok {
		return ErrKeyNotFound
	}

	if m, ok := output.(Unmarshaller); ok {
		return m.UnmarshalSession(value)
	}

	return mapstructure.Decode(value, output)
}

func WithContext(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, sessionContextKey, session)
}

func fromContext(ctx context.Context) (*Session, bool) {
	v := ctx.Value(sessionContextKey)
	if v == nil {
		return nil, false
	}

	session, ok := v.(*Session)
	if !ok {
		return nil, false
	}

	return session, true
}
