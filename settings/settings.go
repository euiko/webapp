package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/mitchellh/mapstructure"
)

type (
	Settings struct {
		Log    Log            `mapstructure:"log"`
		Server Server         `mapstructure:"server"`
		DB     Database       `mapstructure:"db"`
		Extra  map[string]any `mapstructure:"extra"`
	}

	Log struct {
		Level string `mapstructure:"level"`
	}

	Server struct {
		Addr         string        `mapstructure:"addr"`
		ReadTimeout  time.Duration `mapstructure:"read_timeout"`
		WriteTimeout time.Duration `mapstructure:"write_timeout"`
		IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
		ApiPrefix    string        `mapstructure:"api_prefix"`
		// TODO: add https support
	}

	Database struct {
		Sql   SqlDatabase   `mapstructure:"sql"`
		Extra ExtraDatabase `mapstructure:"extra"`
	}

	ExtraDatabase struct {
		Sql map[string]SqlDatabase `mapstructure:"sql"`
	}

	SqlDatabase struct {
		Enabled         bool          `mapstructure:"enabled"`
		Driver          string        `mapstructure:"driver"`
		Uri             string        `mapstructure:"uri"`
		ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
		MaxIdleConns    int           `mapstructure:"max_idle_conns"`
		MaxOpenConns    int           `mapstructure:"max_open_conns"`
		UseORM          bool          `mapstructure:"use_orm"`
	}

	Format int
)

const (
	FormatYaml = iota
	FormatJson
)

var (
	ErrKeyNotFound    = errors.New("key not found")
	ErrMismatchedType = errors.New("type mismatch")
)

func New() Settings {
	return Settings{
		Log: Log{
			Level: "info",
		},
		Server: Server{
			Addr:         ":8080",
			ReadTimeout:  60 * time.Second,
			WriteTimeout: 60 * time.Second,
			IdleTimeout:  0,
			ApiPrefix:    "/api",
			// TODO: add https support
		},
		DB: Database{
			Sql: SqlDatabase{
				Enabled:         false,
				Driver:          "pgx",
				Uri:             "postgres://postgres:12345678@localhost:5432/postgres?sslmode=disable",
				ConnMaxLifetime: 60 * time.Second,
				MaxIdleConns:    10,
				MaxOpenConns:    10,
			},
			Extra: ExtraDatabase{
				Sql: make(map[string]SqlDatabase),
			},
		},
		Extra: make(map[string]any),
	}
}

func GetExtra[T any](s *Settings, key string) (*T, error) {
	extra, ok := s.Extra[key]
	if !ok {
		return nil, ErrKeyNotFound
	}

	switch v := extra.(type) {
	case *T:
		return v, nil
	case T:
		return &v, nil
	default:
		return nil, ErrMismatchedType
	}
}

func Write(s *Settings, format Format, w io.Writer) error {
	var (
		mapCoded = make(map[string]any)
		encoded  []byte
		err      error
	)

	// convert into map
	if err := mapstructure.Decode(s, &mapCoded); err != nil {
		return err
	}

	switch format {
	case FormatYaml:
		encoded, err = yaml.Marshal(mapCoded)
	case FormatJson:
		encoded, err = json.Marshal(mapCoded)
	default:
		return fmt.Errorf("unsupported format: %d", format)
	}

	if err != nil {
		return err
	}

	_, err = w.Write(encoded)
	return err
}
