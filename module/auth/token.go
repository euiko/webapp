package auth

import (
	"context"
	"errors"

	"github.com/euiko/webapp/pkg/token"
	"github.com/euiko/webapp/settings"
	"github.com/lestrrat-go/jwx/v3/jwa"
)

type (
	TokenEncodingFactory func(s *TokenEncodingSettings) token.Encoding

	contextKey struct{}
)

var (
	tokenEncodingRegistry = map[string]TokenEncodingFactory{
		"jwt":          jwtEncodingFactory,
		"headless-jwt": headlessJwtEncodingFactory,
	}

	tokenContextKey = contextKey{}
)

func NewTokenEncoding(s *settings.Settings) (token.Encoding, error) {
	var authSettings Settings
	if err := s.GetExtra("auth", &authSettings); err != nil {
		return nil, err
	}

	fn, ok := tokenEncodingRegistry[authSettings.TokenEncoding.Type]
	if !ok {
		return nil, errors.New("invalid token encoding type (valid types: jwt, headless-jwt)")
	}

	return fn(&authSettings.TokenEncoding), nil
}

func TokenFromContext(ctx context.Context) (*token.Token, bool) {
	v := ctx.Value(tokenContextKey)
	if v == nil {
		return nil, false
	}

	token, ok := v.(*token.Token)
	if !ok {
		return nil, false
	}

	return token, true
}

func IsAuthenticated(ctx context.Context) bool {
	_, ok := TokenFromContext(ctx)
	return ok
}

func contextWithToken(ctx context.Context, token *token.Token) context.Context {
	return context.WithValue(ctx, tokenContextKey, token)
}

func jwtEncodingFactory(s *TokenEncodingSettings) token.Encoding {
	return newJwtEncoding(s, func(sa jwa.SignatureAlgorithm, kp token.KeyProvider, jo ...token.JwtOption) token.Encoding {
		return token.NewJwtEncoding(sa, kp, jo...)
	})
}

func headlessJwtEncodingFactory(s *TokenEncodingSettings) token.Encoding {
	return newJwtEncoding(s, func(sa jwa.SignatureAlgorithm, kp token.KeyProvider, jo ...token.JwtOption) token.Encoding {
		return token.NewHeadlessJwtEncoding(sa, kp, jo...)
	})
}

func newJwtEncoding(
	s *TokenEncodingSettings,
	fn func(jwa.SignatureAlgorithm, token.KeyProvider, ...token.JwtOption) token.Encoding,
) token.Encoding {
	var (
		keyProvider token.KeyProvider
		algorithm   = jwa.NewSignatureAlgorithm(s.JWTAlgorithm)
		hsKeys      = token.NewSymetricKey([]byte(s.HSKey))
	)

	if algorithm.IsSymmetric() {
		keyProvider = hsKeys
	}

	opts := []token.JwtOption{
		token.JwtWithIssuer(s.JWTIssuer),
		token.JwtWithAudience(s.JWTAudience),
		token.JwtWithExpiration(s.JWTTimeout),
	}
	return fn(algorithm, keyProvider, opts...)
}
