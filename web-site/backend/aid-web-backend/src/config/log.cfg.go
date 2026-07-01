package config

type LogCfg struct {
	Daemon DaemonLogCfg `mapstructure:"daemons" json:"daemons" yaml:"daemons" toml:"daemons"`
	Zap    ZapCfg       `mapstructure:"zap" json:"zap" yaml:"zap" toml:"zap"`
}
