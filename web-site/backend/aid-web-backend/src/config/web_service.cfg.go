package config

type (
	WebServiceCfg struct {
		Port       string         `mapstructure:"port" json:"port" yaml:"port" toml:"port"`
		Cors       bool           `mapstructure:"cors" json:"cors" yaml:"cors" toml:"cors"`
		StaticDirs []StaticDirCfg `mapstructure:"static-dir" json:"static-dir" yaml:"static-dir" toml:"static-dir"`
	}

	StaticDirCfg struct {
		URL string `mapstructure:"url" json:"url" yaml:"url" toml:"url"`
		Dir string `mapstructure:"dir" json:"dir" yaml:"dir" toml:"dir"`
	}
)
