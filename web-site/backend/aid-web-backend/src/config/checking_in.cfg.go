package config

// CheckingInCfg 考勤计算临时目录
type CheckingInCfg struct {
	TmpDir string `mapstructure:"tmp-dir" json:"tmp-dir" yaml:"tmp-dir" toml:"tmp-dir"`
}
