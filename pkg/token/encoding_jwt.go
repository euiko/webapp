package token

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

type (
	JwtEncoding struct {
		algorithm   jwa.SignatureAlgorithm
		keyProvider KeyProvider

		useNBF   bool
		now      func() time.Time
		issuer   string
		audience string
		ttl      time.Duration
	}

	HeadlessJwtEncoding struct {
		encoding *JwtEncoding
		header   []byte
	}

	JwtOption func(*JwtEncoding)
)

// JwtWithExpiration sets the expiration time of the token to the current time plus the given duration.
func JwtWithExpiration(ttl time.Duration) JwtOption {
	return func(e *JwtEncoding) {
		e.ttl = ttl
	}
}

// JwtWithIssuer sets the issuer of the token to the given string.
func JwtWithIssuer(issuer string) JwtOption {
	return func(e *JwtEncoding) {
		e.issuer = issuer
	}
}

// WithNotBefore sets the NotBefore claim of the token to the current time.
func JwtWithNotBefore() JwtOption {
	return func(e *JwtEncoding) {
		e.useNBF = true
	}
}

// JwtWithAudience sets the audience of the token to the given string.
func JwtWithAudience(audience string) JwtOption {
	return func(e *JwtEncoding) {
		e.audience = audience
	}
}

// JwtWithTimeProvider sets the time provider for the token to the given function.
func JwtWithTimeProvider(now func() time.Time) JwtOption {
	return func(e *JwtEncoding) {
		e.now = now
	}
}

// NewJwtEncoding creates a new JwtEncoding with the given options.
func NewJwtEncoding(algorithm jwa.SignatureAlgorithm, keyProvider KeyProvider, opts ...JwtOption) *JwtEncoding {
	e := JwtEncoding{
		algorithm:   algorithm,
		keyProvider: keyProvider,
		useNBF:      false,
		now:         time.Now,
		issuer:      "",
		audience:    "",
		ttl:         24 * time.Hour, // default only valid for 24 hour
	}

	for _, opt := range opts {
		opt(&e)
	}

	return &e
}

// NewHeadlessJwtEncoding .
func NewHeadlessJwtEncoding(algorithm jwa.SignatureAlgorithm, keyProvider KeyProvider, opts ...JwtOption) *HeadlessJwtEncoding {
	jwtEncoding := NewJwtEncoding(algorithm, keyProvider, opts...)

	headerData := map[string]interface{}{
		"alg": algorithm.String(),
		"typ": "JWT",
	}
	header, _ := json.Marshal(headerData)

	e := HeadlessJwtEncoding{
		encoding: jwtEncoding,
		header:   header,
	}
	return &e
}

func (e *JwtEncoding) Encode(subject string, audiences ...string) ([]byte, error) {
	now := e.now()
	builder := jwt.NewBuilder().
		Subject(subject).
		Audience(audiences).
		IssuedAt(now).
		Issuer(e.issuer)

	if e.ttl > 0 {
		builder.Expiration(now.Add(e.ttl))
	}

	if e.useNBF {
		builder.NotBefore(now)
	}

	jwtToken, err := builder.Build()
	if err != nil {
		return nil, err
	}

	return jwt.Sign(jwtToken, jwt.WithKey(e.algorithm, e.keyProvider.Private()))
}

func (e *JwtEncoding) Decode(b []byte) (*Token, error) {
	options := []jwt.ParseOption{
		jwt.WithKey(e.algorithm, e.keyProvider.Public()),
	}

	if e.issuer != "" {
		options = append(options, jwt.WithIssuer(e.issuer))
	}

	if e.audience != "" {
		options = append(options, jwt.WithAudience(e.audience))
	}

	verified, err := jwt.Parse(b, options...)
	if err != nil {
		return nil, err
	}

	token := new(Token)
	if issuer, ok := verified.Issuer(); ok {
		token.Issuer = issuer
	}

	if subject, ok := verified.Subject(); ok {
		token.Subject = subject
	}

	if audience, ok := verified.Audience(); ok {
		token.Audience = audience
	}

	if expiresAt, ok := verified.Expiration(); ok {
		token.ExpiresAt = expiresAt
	}

	if issuedAt, ok := verified.IssuedAt(); ok {
		token.IssuedAt = issuedAt
	}

	return token, nil
}

func (e *HeadlessJwtEncoding) Encode(subject string, audiences ...string) ([]byte, error) {
	encoded, err := e.encoding.Encode(subject, audiences...)
	if err != nil {
		return nil, err
	}

	// strip header
	splitted := strings.Split(string(encoded), ".")
	if len(splitted) != 3 {
		return nil, fmt.Errorf("invalid token")
	}
	headless := strings.Join(splitted[1:], ".")

	return []byte(headless), nil
}

func (e *HeadlessJwtEncoding) Decode(b []byte) (*Token, error) {
	// append header first
	b = append(e.header, b...)

	return e.encoding.Decode(b)
}
