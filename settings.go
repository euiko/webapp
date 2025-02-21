package webapp

import (
	"os"
	"path"

	"github.com/euiko/webapp/api"
	"github.com/spf13/viper"
)

func (a *App) loadSettings() *viper.Viper {
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

	return v
}
