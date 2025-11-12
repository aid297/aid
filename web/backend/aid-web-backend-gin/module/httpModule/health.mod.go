package httpModule

import "time"

type (
	HealthResponse struct {
		Time         HealthTime       `json:"time" yaml:"time" toml:"time"`
		System       HealthSystem     `json:"system" yaml:"system" toml:"system"`
		WebService   HealthWebService `json:"webService" yaml:"webService" toml:"webService"`
		VSCodeLaunch HealthVSCode     `json:"vscodeLaunchFile" yaml:"vscodeLaunchFile" toml:"vscodeLaunchFile"`
	}

	HealthTime struct {
		Now    time.Time `json:"now" yaml:"now" toml:"now"`
		String string    `json:"string" yaml:"string" toml:"string"`
	}

	HealthSystem struct {
		Debug    bool   `json:"debug" yaml:"debug" toml:"debug"`
		Version  string `json:"version" yaml:"version" toml:"version"`
		Daemon   bool   `json:"daemon" yaml:"daemon" toml:"daemon"`
		Timezone string `json:"timezone" yaml:"timezone" toml:"timezone"`
	}

	HealthWebService struct {
		Cors bool `json:"cors" yaml:"cors" toml:"cors"`
	}

	HealthVSCode struct {
		Version        string                      `json:"version" yaml:"version" toml:"version"`
		Configurations []HealthVSCodeConfiguration `json:"configurations" yaml:"configurations" toml:"configurations"`
	}

	HealthVSCodeConfiguration struct {
		Name       string            `json:"name" yaml:"name" toml:"name"`
		Type       string            `json:"type" yaml:"type" toml:"type"`
		Request    string            `json:"request" yaml:"request" toml:"request"`
		Mode       string            `json:"mode" yaml:"mode" toml:"mode"`
		Program    string            `json:"program" yaml:"program" toml:"program"`
		Env        map[string]string `json:"env" yaml:"env" toml:"env"`
		Args       []string          `json:"args" yaml:"args" toml:"args"`
		BuildFlags []string          `json:"buildFlags" yaml:"buildFlags" toml:"buildFlags"`
	}
)
