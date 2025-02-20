package settings

import (
	"time"

	"github.com/spf13/viper"
)

type (
	Settings struct {
		Log          Log          `mapstructure:"log"`
		Server       Server       `mapstructure:"server"`
		StaticServer StaticServer `mapstructure:"static_server"`
		DB           Database     `mapstructure:"db"`
		Extra        *viper.Viper `mapstructure:"extra"`
	}

	Log struct {
		Level string `mapstructure:"level"`
	}

	Server struct {
		Addr         string        `mapstructure:"addr"`
		ReadTimeout  time.Duration `mapstructure:"read_timeout"`
		WriteTimeout time.Duration `mapstructure:"write_timeout"`
		IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
		// TODO: add https support
	}

	StaticServer struct {
		Enabled bool        `mapstructure:"enabled"`
		Embed   StaticEmbed `mapstructure:"embed"`
		Proxy   StaticProxy `mapstructure:"proxy"`
	}

	StaticEmbed struct {
		IndexPath string `mapstructure:"index_path"`
		UseMPA    bool   `mapstructure:"use_mpa"`
	}

	StaticProxy struct {
		Upstream string `mapstructure:"upstream"`
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
)

func DefaultSettings() Settings {
	// default settings
	return Settings{
		Log: Log{
			Level: "info",
		},
		Server: Server{
			Addr:         ":8080",
			ReadTimeout:  60 * time.Second,
			WriteTimeout: 60 * time.Second,
			IdleTimeout:  0,
			// TODO: add https support
		},
		StaticServer: StaticServer{
			Enabled: true,
			Embed: StaticEmbed{
				IndexPath: "index.html",
				UseMPA:    false,
			},
			Proxy: StaticProxy{
				Upstream: "http://localhost:5173",
			},
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
		Extra: nil,
	}
}
