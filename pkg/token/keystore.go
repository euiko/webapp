package token

import "sync"

type (
	KeyStore struct {
		keys sync.Map
	}
)

func (ks *KeyStore) Add(id string, key Key) {
	ks.keys.Store(id, key)
}

func (ks *KeyStore) Remove(id string) {
	ks.keys.Delete(id)
}

func (ks *KeyStore) Keys() []Key {
	var keys []Key
	ks.keys.Range(func(_, value any) bool {
		key := value.(Key)
		keys = append(keys, key)
		return true
	})

	return keys
}
