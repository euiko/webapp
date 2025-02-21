package static

type (
	Settings struct {
		Enabled bool          `mapstructure:"enabled"`
		Embed   EmbedSettings `mapstructure:"embed"`
		Proxy   ProxySettings `mapstructure:"proxy"`
	}

	EmbedSettings struct {
		IndexPath string `mapstructure:"index_path"`
		UseMPA    bool   `mapstructure:"use_mpa"`
	}

	ProxySettings struct {
		Upstream string `mapstructure:"upstream"`
	}
)
