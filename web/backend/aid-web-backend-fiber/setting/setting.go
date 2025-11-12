package setting

type Setting struct {
	System     SystemSet     `mapstructure:"system" json:"system" yaml:"system"`
	WebService WebServiceSet `mapstructure:"web-service" json:"web-service" yaml:"web-service"`
	Log        LogSet        `mapstructure:"log" json:"log" yaml:"log"`
}
