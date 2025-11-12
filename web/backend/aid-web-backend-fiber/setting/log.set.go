package setting

type LogSet struct {
	Daemon string `mapstructure:"daemon" json:"daemon" yaml:"daemon"`
	Zap    ZapSet `mapstructure:"zap" json:"zap" yaml:"zap"`
}
