package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/euiko/webapp/pkg/log"
	"github.com/go-viper/mapstructure/v2"
	"github.com/go-yaml/yaml"
)

type (
	Settings struct {
		Log    Log      `mapstructure:"log"`
		Server Server   `mapstructure:"server"`
		DB     Database `mapstructure:"db"`

		extra map[string]any `mapstructure:"extra"`
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
		extra: make(map[string]any),
	}
}

func (s *Settings) SetExtra(key string, value interface{}) {
	s.extra[key] = value
}

func (s *Settings) GetExtra(key string, output interface{}) error {
	outVal := reflect.ValueOf(output)
	if outVal.Kind() != reflect.Ptr {
		return errors.New("result must be a pointer")
	}

	outVal = outVal.Elem()
	if !outVal.CanAddr() {
		return errors.New("result must be addressable pointer")
	}

	value, ok := s.extra[key]
	if !ok {
		return ErrKeyNotFound
	}

	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if reflect.TypeOf(outVal) == reflect.TypeOf(val) {
		outVal.Set(val)
		return nil
	} else if val.Kind() == reflect.Map {
		return mapstructure.Decode(value, outVal.Interface())
	}

	return fmt.Errorf("type mismatch")
}

func Write(s *Settings, format Format, w io.Writer) error {
	var (
		mapCoded = make(map[string]interface{})
		encoded  []byte
		err      error
	)

	// convert into map
	if err := mapstructure.Decode(s, &mapCoded); err != nil {
		return err
	}

	// add extra settings
	extra := make(map[string]interface{})
	hook := mapstructure.DecodeHookFuncValue(structToMapHook)
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: hook,
		Result:     &extra,
	})
	if err != nil {
		return err
	}
	if err := decoder.Decode(s.extra); err != nil {
		return err
	}
	mapCoded["extra"] = extra

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

func structToMapHook(from reflect.Value, to reflect.Value) (interface{}, error) {
	var (
		fromVal = from
		result  = from.Interface()
	)

	log.Info("structToMapHook",
		log.WithField("typeof(from)", reflect.TypeOf(from.Interface())),
		log.WithField("from", from.Interface()),
		log.WithField("to", to.Interface()),
	)

	if from.Kind() == reflect.Ptr {
		fromVal = from.Elem()
	}

	if fromVal.Kind() == reflect.Struct {
		newVal := make(map[string]interface{})
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			Result: &newVal,
		})
		if err != nil {
			return nil, err
		}

		if err := decoder.Decode(fromVal.Interface()); err != nil {
			return nil, err
		}
		result = newVal
	}

	return result, nil
}
