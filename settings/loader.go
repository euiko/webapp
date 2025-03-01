package settings

import (
	"errors"
	"os"
	"path"
	"strings"

	// viper 1.19 still using legacy mapstructure

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type (
	Loader interface {
		Load(*Settings) error
	}

	loader struct {
		v *viper.Viper
	}
)

var (
	ErrConfigNotFound = errors.New("configuration file not found")
)

func NewLoader(name string, envPrefix string) Loader {
	// use viper for configuration
	v := viper.New()
	v.SetConfigName(name)
	v.AddConfigPath(".")

	// use short name as env prefix if it is defined
	if envPrefix != "" {
		v.SetEnvPrefix(envPrefix)
	}

	// add config in home directory is it is defined
	homeDir := os.Getenv("HOME")
	if homeDir != "" {
		v.AddConfigPath(homeDir)
		v.AddConfigPath(path.Join(homeDir, ".config", name))
	}

	loader := loader{
		v: v,
	}
	return &loader
}

func (l *loader) Load(settings *Settings) error {
	// read and load from file
	if err := l.v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return ErrConfigNotFound
		}

		return err
	}

	// decode options
	opts := []viper.DecoderConfigOption{
		viper.DecoderConfigOption(func(dc *mapstructure.DecoderConfig) {
			dc.ErrorUnused = true
		}),
	}

	// unmarshal settings from viper
	err := l.v.Unmarshal(settings, opts...)
	if err != nil {
		// build new error excluding "extra" key
		var (
			newErr  = new(mapstructure.Error)
			origErr = err.(*mapstructure.Error)
		)

		for _, err := range origErr.Errors {
			if strings.Contains(err, "has invalid keys: extra") {
				continue
			}
			newErr.Errors = append(newErr.Errors, err)
		}

		// only return error when the newErr is not empty
		if len(newErr.Errors) > 0 {
			return newErr
		}
	}

	// unmarshal extra settings from viper
	for key, value := range settings.extra {
		if err := l.v.UnmarshalKey("extra."+key, &value, opts...); err != nil {
			return err
		}
	}

	return nil
}
