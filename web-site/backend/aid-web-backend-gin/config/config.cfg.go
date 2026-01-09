package config

// Config 总配置
type Config struct {
	System     SystemCfg     `mapstructure:"system" json:"system" yaml:"system" toml:"system"`
	WebService WebServiceCfg `mapstructure:"web-service" json:"web-service" yaml:"web-service" toml:"web-service"`
	Log        LogCfg        `mapstructure:"log" json:"log" yaml:"log" toml:"log"`
	Rezip      RezipCfg      `mapstructure:"rezip" json:"rezip" yaml:"rezip" toml:"rezip"`
	UploadDir  string        `mapstructure:"upload-dir" json:"upload-dir" yaml:"upload-dir" toml:"upload-dir"`
}
