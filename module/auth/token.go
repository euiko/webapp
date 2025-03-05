package auth

import (
	"errors"
	"fmt"

	"github.com/euiko/webapp/pkg/token"
	"github.com/lestrrat-go/jwx/v3/jwa"
)

type (
	TokenEncodingFactory func(s *TokenEncodingSettings) (token.Encoding, error)
)

var (
	tokenEncodingRegistry = map[string]TokenEncodingFactory{
		"jwt":          jwtEncodingFactory,
		"headless-jwt": headlessJwtEncodingFactory,
	}
)

func NewTokenEncoding(s *Settings) (token.Encoding, error) {
	fn, ok := tokenEncodingRegistry[s.TokenEncoding.Type]
	if !ok {
		return nil, errors.New("invalid token encoding type (valid types: jwt, headless-jwt)")
	}

	return fn(&s.TokenEncoding)
}

func jwtEncodingFactory(s *TokenEncodingSettings) (token.Encoding, error) {
	return newJwtEncoding(s, func(sa jwa.SignatureAlgorithm, kp token.KeyProvider, jo ...token.JwtOption) token.Encoding {
		return token.NewJwtEncoding(sa, kp, jo...)
	})
}

func headlessJwtEncodingFactory(s *TokenEncodingSettings) (token.Encoding, error) {
	return newJwtEncoding(s, func(sa jwa.SignatureAlgorithm, kp token.KeyProvider, jo ...token.JwtOption) token.Encoding {
		return token.NewHeadlessJwtEncoding(sa, kp, jo...)
	})
}

func newJwtEncoding(
	s *TokenEncodingSettings,
	fn func(jwa.SignatureAlgorithm, token.KeyProvider, ...token.JwtOption) token.Encoding,
) (token.Encoding, error) {
	var (
		keyProvider token.KeyProvider
		hsKeys      = token.NewSymetricKey([]byte(s.HSKey))
	)

	algorithm, ok := jwa.LookupSignatureAlgorithm(s.JWTAlgorithm)
	if !ok {
		return nil, fmt.Errorf("invalid jwt algorithm: %s", s.JWTAlgorithm)
	}

	if algorithm.IsSymmetric() {
		keyProvider = hsKeys
	}

	opts := []token.JwtOption{
		token.JwtWithIssuer(s.JWTIssuer),
		token.JwtWithAudience(s.JWTAudience),
		token.JwtWithExpiration(s.JWTTimeout),
	}
	return fn(algorithm, keyProvider, opts...), nil
}
