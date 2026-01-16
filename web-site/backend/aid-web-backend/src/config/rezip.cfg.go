package config

// FileManagerCfg 文件管理配置
type FileManagerCfg struct {
	Dir string `mapstructure:"dir" json:"dir" yaml:"dir" toml:"dir"`
}
