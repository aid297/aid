package config

// RezipCfg 重复压缩配置
type RezipCfg struct {
	OutDir string `mapstructure:"out-dir" json:"out-dir" yaml:"out-dir" toml:"out-dir"`
}
