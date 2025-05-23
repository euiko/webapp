package token

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

type (
	JwtEncoding struct {
		algorithm jwa.SignatureAlgorithm

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
func NewJwtEncoding(algorithm jwa.SignatureAlgorithm, opts ...JwtOption) *JwtEncoding {
	e := JwtEncoding{
		algorithm: algorithm,
		useNBF:    false,
		now:       time.Now,
		issuer:    "",
		audience:  "",
		ttl:       24 * time.Hour, // default only valid for 24 hour
	}

	for _, opt := range opts {
		opt(&e)
	}

	return &e
}

// NewHeadlessJwtEncoding .
func NewHeadlessJwtEncoding(algorithm jwa.SignatureAlgorithm, opts ...JwtOption) *HeadlessJwtEncoding {
	jwtEncoding := NewJwtEncoding(algorithm, opts...)

	headerData := map[string]interface{}{
		"alg": algorithm.String(),
		"typ": "JWT",
	}
	header, _ := json.Marshal(headerData)
	base64Header := base64.StdEncoding.EncodeToString(header)

	e := HeadlessJwtEncoding{
		encoding: jwtEncoding,
		header:   []byte(base64Header),
	}
	return &e
}

func (e *JwtEncoding) Encode(key Key, subject string, audiences ...string) ([]byte, error) {
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

	return jwt.Sign(jwtToken, jwt.WithKey(e.algorithm, key.Private()))
}

func (e *JwtEncoding) Decode(key Key, b []byte) (*Token, error) {
	options := []jwt.ParseOption{
		jwt.WithKey(e.algorithm, key.Public()),
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

func (e *HeadlessJwtEncoding) Encode(key Key, subject string, audiences ...string) ([]byte, error) {
	encoded, err := e.encoding.Encode(key, subject, audiences...)
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

func (e *HeadlessJwtEncoding) Decode(key Key, b []byte) (*Token, error) {
	// append header first
	headerBytes := append(e.header, '.')
	b = append(headerBytes, b...)
	return e.encoding.Decode(key, b)
}
