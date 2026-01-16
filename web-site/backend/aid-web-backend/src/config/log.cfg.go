package config

type LogCfg struct {
	Daemon DaemonLogCfg `mapstructure:"daemon" json:"daemon" yaml:"daemon" toml:"daemon"`
	Zap    ZapCfg       `mapstructure:"zap" json:"zap" yaml:"zap" toml:"zap"`
}
