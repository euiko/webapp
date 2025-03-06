package helper

import (
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"hash"
)

type (
	HashAlgorithm int
)

const (
	HashSHA256 HashAlgorithm = iota
	HashSHA512
)

var (
	ErrInvalidHashAlgorithm = errors.New("invalid hash algorithm")

	hashAlgorithms = map[HashAlgorithm]func() hash.Hash{
		HashSHA256: sha256.New,
		HashSHA512: sha512.New,
	}
)

func Hash(data []byte, hashType HashAlgorithm) []byte {
	factory := hashAlgorithms[hashType]
	hasher := factory()
	return hasher.Sum(data)
}
