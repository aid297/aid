package config

// FileManagerCfg 文件管理配置
type FileManagerCfg struct {
	Port string `mapstructure:"port" json:"port" yaml:"port" toml:"port"`
	Dir  string `mapstructure:"dir" json:"dir" yaml:"dir" toml:"dir"`
}
