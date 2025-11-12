package config

type DaemonLogCfg struct {
	Dir      string `mapstructure:"dir" json:"dir" yaml:"dir" toml:"dir"`
	Filename string `mapstructure:"filename" json:"filename" yaml:"filename" toml:"filename"`
}
