package webapp

import (
	"os"
	"path"

	"github.com/euiko/webapp/settings"
	"github.com/spf13/viper"
)

func loadSettings(name string, shortName string) settings.Settings {
	s := settings.DefaultSettings()

	// use viper for configuration
	v := viper.New()
	v.SetConfigName(name)
	v.AddConfigPath(".")

	// use short name as env prefix if it is defined
	if shortName != "" {
		v.SetEnvPrefix(shortName)
	}

	// add config in home directory is it is defined
	homeDir := os.Getenv("HOME")
	if homeDir != "" {
		v.AddConfigPath(homeDir)
		v.AddConfigPath(path.Join(homeDir, ".config", name))
	}

	// load settings
	v.Unmarshal(&s)
	return s
}
