package auth

import "time"

type (
	Settings struct {
		Enabled       bool                  `mapstructure:"enabled"`
		TokenEncoding TokenEncodingSettings `mapstructure:"token_encoding"`
	}

	TokenEncodingSettings struct {
		Type         string        `mapstructure:"type"`
		JWTAlgorithm string        `mapstructure:"jwt_algorithm"`
		JWTIssuer    string        `mapstructure:"jwt_issuer"`
		JWTAudience  string        `mapstructure:"jwt_audience"`
		JWTTimeout   time.Duration `mapstructure:"jwt_timeout"`
		Keys         []string      `mapstructure:"keys"`
	}
)
