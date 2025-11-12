package setting

type ZapSet struct {
	Level       string `mapstructure:"level" json:"level" yaml:"level"`
	EncoderType string `mapstructure:"encoder-type" json:"encoder-type" yaml:"encoder-type"`
	Extension   string `mapstructure:"extension" json:"extension" yaml:"extension"`
	InConsole   bool   `mapstructure:"in-console" json:"in-console" yaml:"in-console"`
	MaxSize     int    `mapstructure:"max-size" json:"max-size" yaml:"max-size"`
	MaxDay      int    `mapstructure:"max-day" json:"max-day" yaml:"max-day"`
	DirAbs      bool   `mapstructure:"dir-abs" json:"dir-abs" yaml:"dir-abs"`
	Dir         string `mapstructure:"dir" json:"dir" yaml:"dir"`
}
