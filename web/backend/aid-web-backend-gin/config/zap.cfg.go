package config

type ZapCfg struct {
	Level       string `mapstructure:"level" json:"level" yaml:"level" toml:"level"`
	EncoderType string `mapstructure:"encoder-type" json:"encoder-type" yaml:"encoder-type" toml:"encoder-type"`
	Extension   string `mapstructure:"extension" json:"extension" yaml:"extension" toml:"extension"`
	InConsole   bool   `mapstructure:"in-console" json:"in-console" yaml:"in-console" toml:"in-console"`
	MaxSize     int    `mapstructure:"max-size" json:"max-size" yaml:"max-size" toml:"max-size"`
	MaxDay      int    `mapstructure:"max-day" json:"max-day" yaml:"max-day" toml:"max-day"`
	DirAbs      bool   `mapstructure:"dir-abs" json:"dir-abs" yaml:"dir-abs" toml:"dir-abs"`
	Dir         string `mapstructure:"dir" json:"dir" yaml:"dir" toml:"dir"`
}
