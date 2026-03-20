package config

// MessageBoardCfg 留言板配置
type MessageBoardCfg struct {
	Dir string `mapstructure:"dir" json:"dir" yaml:"dir" toml:"dir"`
}
