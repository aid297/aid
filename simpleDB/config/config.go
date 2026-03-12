package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const DefaultConfigPath = "simpleDB/config.yaml"

type (
	Config struct {
		Database  DatabaseConfig  `yaml:"database"`
		Transport TransportConfig `yaml:"transport"`
	}

	DatabaseConfig struct {
		Path string `yaml:"path"`
	}

	TransportConfig struct {
		HTTP HTTPConfig `yaml:"http"`
	}

	HTTPRouteConfig struct {
		Health               string `yaml:"health"`
		Profile              string `yaml:"profile"`
		Login                string `yaml:"login"`
		Register             string `yaml:"register"`
		Refresh              string `yaml:"refresh"`
		Logout               string `yaml:"logout"`
		Activate             string `yaml:"activate"`
		Deactivate           string `yaml:"deactivate"`
		AssignRole           string `yaml:"assignRole"`
		AssignRolePermission string `yaml:"assignRolePermission"`
		InitSDBPassword      string `yaml:"initSDBPassword"`
		SQLExecute           string `yaml:"sql"`
		Admin                string `yaml:"admin"`
		Report               string `yaml:"report"`
	}

	HTTPConfig struct {
		Enabled      bool            `yaml:"enabled"`
		Address      string          `yaml:"address"`
		GinMode      string          `yaml:"ginMode"`
		Route        HTTPRouteConfig `yaml:"route"`
		InitPassword string          `yaml:"initPassword"`
		TokenTTL     string          `yaml:"tokenTTL"`
		TokenSecret  string          `yaml:"tokenSecret"`
		EnableAdmin  bool            `yaml:"enableAdminRoute"`
		EnableReport bool            `yaml:"enableReportRoute"`
		// SQLAllowedOps 控制 /sql/execute 接口允许执行的语句类型。
		// 可选值： select, insert, update, delete, create, drop, truncate, alter
		// 留空表示不限制（默认开放全部，建议生产环境按需收紧）。
		SQLAllowedOps []string `yaml:"sqlAllowedOps"`
	}
)

func Default() Config {
	return Config{
		Database: DatabaseConfig{Path: "demo_transport"},
		Transport: TransportConfig{HTTP: HTTPConfig{
			Enabled: true,
			Address: ":18080",
			GinMode: "release",
			Route: HTTPRouteConfig{
				Health:               "/health",
				Profile:              "/me",
				Login:                "/auth/login",
				Register:             "/auth/register",
				Refresh:              "/auth/refresh",
				Logout:               "/auth/logout",
				Activate:             "/auth/activate",
				Deactivate:           "/auth/deactivate",
				AssignRole:           "/auth/assign-roles",
				AssignRolePermission: "/auth/assign-role-permissions",
				InitSDBPassword:      "/auth/init-sdb-password",
				SQLExecute:           "/sql/execute",
				Admin:                "/admin",
				Report:               "/reports",
			},
			InitPassword: "",
			TokenTTL:     "12h",
			TokenSecret:  "",
			EnableAdmin:  true,
			EnableReport: true,
		}},
	}
}

func Load(path string) (Config, error) {
	config := Default()
	path = strings.TrimSpace(path)
	if path == "" {
		return config, nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil
		}
		return Config{}, err
	}

	var raw rawConfig
	if err = yaml.Unmarshal(content, &raw); err != nil {
		return Config{}, err
	}
	applyRawConfig(&config, raw)
	config.ApplyDefaults()
	return config, nil
}

func Save(path string, config Config) error {
	config.ApplyDefaults()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	payload, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(path, payload, 0o644)
}

func (c *Config) ApplyDefaults() {
	defaults := Default()
	if strings.TrimSpace(c.Database.Path) == "" {
		c.Database.Path = defaults.Database.Path
	}
	if strings.TrimSpace(c.Transport.HTTP.Address) == "" && !c.Transport.HTTP.Enabled && !c.Transport.HTTP.EnableAdmin && !c.Transport.HTTP.EnableReport && strings.TrimSpace(c.Transport.HTTP.TokenTTL) == "" {
		c.Transport.HTTP.Enabled = defaults.Transport.HTTP.Enabled
		c.Transport.HTTP.EnableAdmin = defaults.Transport.HTTP.EnableAdmin
		c.Transport.HTTP.EnableReport = defaults.Transport.HTTP.EnableReport
	}
	if strings.TrimSpace(c.Transport.HTTP.Address) == "" {
		c.Transport.HTTP.Address = defaults.Transport.HTTP.Address
	}
	if strings.TrimSpace(c.Transport.HTTP.GinMode) == "" {
		c.Transport.HTTP.GinMode = defaults.Transport.HTTP.GinMode
	}
	if strings.TrimSpace(c.Transport.HTTP.Route.Health) == "" {
		c.Transport.HTTP.Route.Health = defaults.Transport.HTTP.Route.Health
	}
	if strings.TrimSpace(c.Transport.HTTP.Route.Profile) == "" {
		c.Transport.HTTP.Route.Profile = defaults.Transport.HTTP.Route.Profile
	}
	if strings.TrimSpace(c.Transport.HTTP.Route.Login) == "" {
		c.Transport.HTTP.Route.Login = defaults.Transport.HTTP.Route.Login
	}
	if strings.TrimSpace(c.Transport.HTTP.Route.Refresh) == "" {
		c.Transport.HTTP.Route.Refresh = defaults.Transport.HTTP.Route.Refresh
	}
	if strings.TrimSpace(c.Transport.HTTP.Route.Register) == "" {
		c.Transport.HTTP.Route.Register = defaults.Transport.HTTP.Route.Register
	}
	if strings.TrimSpace(c.Transport.HTTP.Route.Logout) == "" {
		c.Transport.HTTP.Route.Logout = defaults.Transport.HTTP.Route.Logout
	}
	if strings.TrimSpace(c.Transport.HTTP.Route.Activate) == "" {
		c.Transport.HTTP.Route.Activate = defaults.Transport.HTTP.Route.Activate
	}
	if strings.TrimSpace(c.Transport.HTTP.Route.Deactivate) == "" {
		c.Transport.HTTP.Route.Deactivate = defaults.Transport.HTTP.Route.Deactivate
	}
	if strings.TrimSpace(c.Transport.HTTP.Route.AssignRole) == "" {
		c.Transport.HTTP.Route.AssignRole = defaults.Transport.HTTP.Route.AssignRole
	}
	if strings.TrimSpace(c.Transport.HTTP.Route.AssignRolePermission) == "" {
		c.Transport.HTTP.Route.AssignRolePermission = defaults.Transport.HTTP.Route.AssignRolePermission
	}
	if strings.TrimSpace(c.Transport.HTTP.Route.InitSDBPassword) == "" {
		c.Transport.HTTP.Route.InitSDBPassword = defaults.Transport.HTTP.Route.InitSDBPassword
	}
	if strings.TrimSpace(c.Transport.HTTP.Route.SQLExecute) == "" {
		c.Transport.HTTP.Route.SQLExecute = defaults.Transport.HTTP.Route.SQLExecute
	}
	if strings.TrimSpace(c.Transport.HTTP.TokenTTL) == "" {
		c.Transport.HTTP.TokenTTL = defaults.Transport.HTTP.TokenTTL
	}
	if strings.TrimSpace(c.Transport.HTTP.Route.Admin) == "" {
		c.Transport.HTTP.Route.Admin = defaults.Transport.HTTP.Route.Admin
	}
	if strings.TrimSpace(c.Transport.HTTP.Route.Report) == "" {
		c.Transport.HTTP.Route.Report = defaults.Transport.HTTP.Route.Report
	}
}

func (c Config) ParseTokenTTL() (time.Duration, error) {
	value := strings.TrimSpace(c.Transport.HTTP.TokenTTL)
	if value == "" {
		return 12 * time.Hour, nil
	}
	return time.ParseDuration(value)
}

func (c Config) Validate() error {
	c.ApplyDefaults()
	if strings.TrimSpace(c.Database.Path) == "" {
		return fmt.Errorf("database.path is required")
	}
	if !c.Transport.HTTP.Enabled {
		return fmt.Errorf("transport.http.enabled must be true currently")
	}
	if _, err := c.ParseTokenTTL(); err != nil {
		return fmt.Errorf("transport.http.tokenTTL is invalid: %w", err)
	}
	return nil
}

type rawConfig struct {
	Database  rawDatabaseConfig  `yaml:"database"`
	Transport rawTransportConfig `yaml:"transport"`
}

type rawDatabaseConfig struct {
	Path string `yaml:"path"`
}

type rawTransportConfig struct {
	HTTP rawHTTPConfig `yaml:"http"`
}

type rawHTTPConfig struct {
	Enabled       *bool              `yaml:"enabled"`
	Address       string             `yaml:"address"`
	GinMode       string             `yaml:"ginMode"`
	Route         rawHTTPRouteConfig `yaml:"route"`
	InitPassword  string             `yaml:"initPassword"`
	TokenTTL      string             `yaml:"tokenTTL"`
	TokenSecret   string             `yaml:"tokenSecret"`
	EnableAdmin   *bool              `yaml:"enableAdminRoute"`
	EnableReport  *bool              `yaml:"enableReportRoute"`
	SQLAllowedOps []string           `yaml:"sqlAllowedOps"`

	// legacy flat route keys (for backward compatibility)
	HealthPath               string `yaml:"healthPath"`
	ProfilePath              string `yaml:"profilePath"`
	LoginPath                string `yaml:"loginPath"`
	RegisterPath             string `yaml:"registerPath"`
	RefreshPath              string `yaml:"refreshPath"`
	LogoutPath               string `yaml:"logoutPath"`
	ActivatePath             string `yaml:"activatePath"`
	DeactivatePath           string `yaml:"deactivatePath"`
	AssignRolePath           string `yaml:"assignRolePath"`
	AssignRolePermissionPath string `yaml:"assignRolePermissionPath"`
	InitSDBPasswordPath      string `yaml:"initSDBPasswordPath"`
	SQLExecutePath           string `yaml:"sqlExecutePath"`
	AdminPath                string `yaml:"adminPath"`
	ReportPath               string `yaml:"reportPath"`
}

type rawHTTPRouteConfig struct {
	Health               string `yaml:"health"`
	Profile              string `yaml:"profile"`
	Login                string `yaml:"login"`
	Register             string `yaml:"register"`
	Refresh              string `yaml:"refresh"`
	Logout               string `yaml:"logout"`
	Activate             string `yaml:"activate"`
	Deactivate           string `yaml:"deactivate"`
	AssignRole           string `yaml:"assignRole"`
	AssignRolePermission string `yaml:"assignRolePermission"`
	InitSDBPassword      string `yaml:"initSDBPassword"`
	SQLExecute           string `yaml:"sql"`
	Admin                string `yaml:"admin"`
	Report               string `yaml:"report"`
}

func applyRawConfig(config *Config, raw rawConfig) {
	if strings.TrimSpace(raw.Database.Path) != "" {
		config.Database.Path = strings.TrimSpace(raw.Database.Path)
	}
	if raw.Transport.HTTP.Enabled != nil {
		config.Transport.HTTP.Enabled = *raw.Transport.HTTP.Enabled
	}
	if strings.TrimSpace(raw.Transport.HTTP.Address) != "" {
		config.Transport.HTTP.Address = strings.TrimSpace(raw.Transport.HTTP.Address)
	}
	if strings.TrimSpace(raw.Transport.HTTP.GinMode) != "" {
		config.Transport.HTTP.GinMode = strings.TrimSpace(raw.Transport.HTTP.GinMode)
	}
	if value := firstNonEmpty(raw.Transport.HTTP.Route.Health, raw.Transport.HTTP.HealthPath); value != "" {
		config.Transport.HTTP.Route.Health = value
	}
	if value := firstNonEmpty(raw.Transport.HTTP.Route.Profile, raw.Transport.HTTP.ProfilePath); value != "" {
		config.Transport.HTTP.Route.Profile = value
	}
	if value := firstNonEmpty(raw.Transport.HTTP.Route.Login, raw.Transport.HTTP.LoginPath); value != "" {
		config.Transport.HTTP.Route.Login = value
	}
	if value := firstNonEmpty(raw.Transport.HTTP.Route.Refresh, raw.Transport.HTTP.RefreshPath); value != "" {
		config.Transport.HTTP.Route.Refresh = value
	}
	if value := firstNonEmpty(raw.Transport.HTTP.Route.Register, raw.Transport.HTTP.RegisterPath); value != "" {
		config.Transport.HTTP.Route.Register = value
	}
	if value := firstNonEmpty(raw.Transport.HTTP.Route.Logout, raw.Transport.HTTP.LogoutPath); value != "" {
		config.Transport.HTTP.Route.Logout = value
	}
	if value := firstNonEmpty(raw.Transport.HTTP.Route.Activate, raw.Transport.HTTP.ActivatePath); value != "" {
		config.Transport.HTTP.Route.Activate = value
	}
	if value := firstNonEmpty(raw.Transport.HTTP.Route.Deactivate, raw.Transport.HTTP.DeactivatePath); value != "" {
		config.Transport.HTTP.Route.Deactivate = value
	}
	if value := firstNonEmpty(raw.Transport.HTTP.Route.AssignRole, raw.Transport.HTTP.AssignRolePath); value != "" {
		config.Transport.HTTP.Route.AssignRole = value
	}
	if value := firstNonEmpty(raw.Transport.HTTP.Route.AssignRolePermission, raw.Transport.HTTP.AssignRolePermissionPath); value != "" {
		config.Transport.HTTP.Route.AssignRolePermission = value
	}
	if value := firstNonEmpty(raw.Transport.HTTP.Route.InitSDBPassword, raw.Transport.HTTP.InitSDBPasswordPath); value != "" {
		config.Transport.HTTP.Route.InitSDBPassword = value
	}
	if value := firstNonEmpty(raw.Transport.HTTP.Route.SQLExecute, raw.Transport.HTTP.SQLExecutePath); value != "" {
		config.Transport.HTTP.Route.SQLExecute = value
	}
	if raw.Transport.HTTP.InitPassword != "" {
		config.Transport.HTTP.InitPassword = raw.Transport.HTTP.InitPassword
	}
	if strings.TrimSpace(raw.Transport.HTTP.TokenTTL) != "" {
		config.Transport.HTTP.TokenTTL = strings.TrimSpace(raw.Transport.HTTP.TokenTTL)
	}
	config.Transport.HTTP.TokenSecret = raw.Transport.HTTP.TokenSecret
	if raw.Transport.HTTP.EnableAdmin != nil {
		config.Transport.HTTP.EnableAdmin = *raw.Transport.HTTP.EnableAdmin
	}
	if raw.Transport.HTTP.EnableReport != nil {
		config.Transport.HTTP.EnableReport = *raw.Transport.HTTP.EnableReport
	}
	if value := firstNonEmpty(raw.Transport.HTTP.Route.Admin, raw.Transport.HTTP.AdminPath); value != "" {
		config.Transport.HTTP.Route.Admin = value
	}
	if value := firstNonEmpty(raw.Transport.HTTP.Route.Report, raw.Transport.HTTP.ReportPath); value != "" {
		config.Transport.HTTP.Route.Report = value
	}
	if len(raw.Transport.HTTP.SQLAllowedOps) > 0 {
		config.Transport.HTTP.SQLAllowedOps = raw.Transport.HTTP.SQLAllowedOps
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}
