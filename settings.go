package webapp

import (
	"errors"
	"os"
	"path"

	"github.com/euiko/webapp/api"
	"github.com/euiko/webapp/pkg/log"
	"github.com/spf13/viper"
)

var ErrConfigNotFound = errors.New("config not found")

func (a *App) loadSettings() (Unmarshaler, error) {
	// use viper for configuration
	v := viper.New()
	v.SetConfigName(a.name)
	v.AddConfigPath(".")

	// use short name as env prefix if it is defined
	if a.shortName != "" {
		v.SetEnvPrefix(a.shortName)
	}

	// add config in home directory is it is defined
	homeDir := os.Getenv("HOME")
	if homeDir != "" {
		v.AddConfigPath(homeDir)
		v.AddConfigPath(path.Join(homeDir, ".config", a.name))
	}

	// settings loader to configure default settings
	for _, module := range a.modules {
		if loader, ok := module.(api.SettingsLoader); ok {
			loader.DefaultSettings(&a.settings)
		}
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Info("configuration file not found, use default settings",
				log.WithError(err),
			)
			return nil, ErrConfigNotFound
		}

		return nil, err
	}

	return UnmarshalerFunc(func(dst any) error {
		return v.Unmarshal(dst)
	}), nil
}
