package setting

type WebServiceSet struct {
	Port string `mapstructure:"port" json:"port" yaml:"port"`
	Cors bool   `mapstructure:"cors" json:"cors" yaml:"cors"`
}
