package setting

type SystemSet struct {
	Debug   bool   `mapstructure:"debug" json:"debug" yaml:"debug"`
	Version string `mapstructure:"version" json:"version" yaml:"version"`
	Daemon  bool   `mapstructure:"daemon" json:"daemon" yaml:"daemon"`
}
