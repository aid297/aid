package config

// SystemCfg 系统配置
type SystemCfg struct {
	Debug    bool   `mapstructure:"debug" json:"debug" yaml:"debug" toml:"debug"`
	Version  string `mapstructure:"version" json:"version" yaml:"version" toml:"version"`
	Daemon   bool   `mapstructure:"daemon" json:"daemon" yaml:"daemon" toml:"daemon"`
	Timezone string `mapstructure:"timezone" json:"timezone" yaml:"timezone" toml:"timezone"`
}
