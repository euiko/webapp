package helper

import (
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"hash"
	"hash/crc32"
	"hash/crc64"
	"hash/fnv"
)

type (
	HashAlgorithm   int
	Hash64Algorithm int
	Hash32Algorithm int
)

const (
	HashSHA256 HashAlgorithm = iota
	HashSHA512
)

const (
	Hash32FNV1 Hash32Algorithm = iota
	Hash32FNV1a
	Hash32CRC
)

const (
	Hash64FNV1 Hash64Algorithm = iota
	Hash64FNV1a
	Hash64CRC
)

var (
	ErrInvalidHashAlgorithm = errors.New("invalid hash algorithm")

	hashAlgorithms = map[HashAlgorithm]func() hash.Hash{
		HashSHA256: sha256.New,
		HashSHA512: sha512.New,
	}

	hash32Algorithms = map[Hash32Algorithm]func() hash.Hash32{
		Hash32FNV1:  fnv.New32,
		Hash32FNV1a: fnv.New32a,
		Hash32CRC: func() hash.Hash32 {
			return crc32.New(crc32.IEEETable)
		},
	}

	hash64Algorithms = map[Hash64Algorithm]func() hash.Hash64{
		Hash64FNV1:  fnv.New64,
		Hash64FNV1a: fnv.New64a,
		Hash64CRC: func() hash.Hash64 {
			return crc64.New(crc64.MakeTable(crc64.ISO))
		},
	}
)

func Hash(data []byte, hashTypes ...HashAlgorithm) []byte {
	hashType := HashSHA256
	if len(hashTypes) > 0 {
		hashType = hashTypes[0]
	}

	factory := hashAlgorithms[hashType]
	hasher := factory()
	return hasher.Sum(data)
}

func Hash32(data []byte, hashTypes ...Hash32Algorithm) uint32 {
	hashType := Hash32FNV1
	if len(hashTypes) > 0 {
		hashType = hashTypes[0]
	}

	factory := hash32Algorithms[hashType]
	hasher := factory()
	hasher.Write(data)
	return hasher.Sum32()
}

func Hash64(data []byte, hashTypes ...Hash64Algorithm) uint64 {
	hashType := Hash64FNV1
	if len(hashTypes) > 0 {
		hashType = hashTypes[0]
	}

	factory := hash64Algorithms[hashType]
	hasher := factory()
	hasher.Write(data)
	return hasher.Sum64()
}
