package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const DefaultConfigPath = "simpleDB/config.yaml"

type (
	Config struct {
		Database  DatabaseConfig  `yaml:"database"`
		Engine    EngineConfig    `yaml:"engine"`
		Transport TransportConfig `yaml:"transport"`
	}

	DatabaseConfig struct {
		Path string `yaml:"path"`
	}

	EngineConfig struct {
		Persistence EnginePersistenceConfig `yaml:"persistence"`
		Security    EngineSecurityConfig    `yaml:"security"`
	}

	EnginePersistenceConfig struct {
		WindowSeconds int    `yaml:"windowSeconds"`
		WindowBytes   string `yaml:"windowBytes"`
		Threshold     string `yaml:"threshold"`
	}

	EngineSecurityConfig struct {
		CompressAlgorithm string `yaml:"compressAlgorithm"`
		EncryptAlgorithm  string `yaml:"encryptAlgorithm"`
		EncryptKey        string `yaml:"encryptKey"`
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
		SQLGrant             string `yaml:"sqlGrant"`
		SQLRevoke            string `yaml:"sqlRevoke"`
		Admin                string `yaml:"admin"`
		Report               string `yaml:"report"`
	}

	HTTPConfig struct {
		Enabled      bool            `yaml:"enabled"`
		Address      string          `yaml:"address"`
		GinMode      string          `yaml:"ginMode"`
		Route        HTTPRouteConfig `yaml:"route"`
		Limit        HTTPLimitConfig `yaml:"limit"`
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

	HTTPLimitConfig struct {
		Enabled      bool     `yaml:"enabled"`
		Requests     int      `yaml:"requests"`
		Window       string   `yaml:"window"`
		NoTokenPaths []string `yaml:"noTokenPaths"`
	}
)

func Default() Config {
	return Config{
		Database: DatabaseConfig{Path: "demo_transport"},
		Engine: EngineConfig{
			Persistence: EnginePersistenceConfig{
				WindowSeconds: 10,
				WindowBytes:   "10mb",
				Threshold:     "100mb",
			},
			Security: EngineSecurityConfig{
				CompressAlgorithm: "",
				EncryptAlgorithm:  "",
				EncryptKey:        "",
			},
		},
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
				SQLGrant:             "/sql/grant",
				SQLRevoke:            "/sql/revoke",
				Admin:                "/admin",
				Report:               "/reports",
			},
			Limit: HTTPLimitConfig{
				Enabled:      false,
				Requests:     60,
				Window:       "1m",
				NoTokenPaths: []string{"/auth/login", "/auth/register", "/auth/init-sdb-password", "/health"},
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
	if c.Engine.Persistence.WindowSeconds <= 0 {
		c.Engine.Persistence.WindowSeconds = defaults.Engine.Persistence.WindowSeconds
	}
	if strings.TrimSpace(c.Engine.Persistence.WindowBytes) == "" {
		c.Engine.Persistence.WindowBytes = defaults.Engine.Persistence.WindowBytes
	}
	if strings.TrimSpace(c.Engine.Persistence.Threshold) == "" {
		c.Engine.Persistence.Threshold = defaults.Engine.Persistence.Threshold
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
	if strings.TrimSpace(c.Transport.HTTP.Route.SQLGrant) == "" {
		c.Transport.HTTP.Route.SQLGrant = defaults.Transport.HTTP.Route.SQLGrant
	}
	if strings.TrimSpace(c.Transport.HTTP.Route.SQLRevoke) == "" {
		c.Transport.HTTP.Route.SQLRevoke = defaults.Transport.HTTP.Route.SQLRevoke
	}
	if strings.TrimSpace(c.Transport.HTTP.TokenTTL) == "" {
		c.Transport.HTTP.TokenTTL = defaults.Transport.HTTP.TokenTTL
	}
	if c.Transport.HTTP.Limit.Requests <= 0 {
		c.Transport.HTTP.Limit.Requests = defaults.Transport.HTTP.Limit.Requests
	}
	if strings.TrimSpace(c.Transport.HTTP.Limit.Window) == "" {
		c.Transport.HTTP.Limit.Window = defaults.Transport.HTTP.Limit.Window
	}
	if len(c.Transport.HTTP.Limit.NoTokenPaths) == 0 {
		c.Transport.HTTP.Limit.NoTokenPaths = append([]string(nil), defaults.Transport.HTTP.Limit.NoTokenPaths...)
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
	if _, err := c.ParseEngineWindowBytes(); err != nil {
		return fmt.Errorf("engine.persistence.windowBytes is invalid: %w", err)
	}
	if _, err := c.ParseEngineThresholdBytes(); err != nil {
		return fmt.Errorf("engine.persistence.threshold is invalid: %w", err)
	}
	if !c.Transport.HTTP.Enabled {
		return fmt.Errorf("transport.http.enabled must be true currently")
	}
	if _, err := c.ParseTokenTTL(); err != nil {
		return fmt.Errorf("transport.http.tokenTTL is invalid: %w", err)
	}
	if c.Transport.HTTP.Limit.Enabled {
		if c.Transport.HTTP.Limit.Requests <= 0 {
			return fmt.Errorf("transport.http.limit.requests must be > 0")
		}
		if _, err := c.ParseLimitWindow(); err != nil {
			return fmt.Errorf("transport.http.limit.window is invalid: %w", err)
		}
	}
	return nil
}

func (c Config) ParseLimitWindow() (time.Duration, error) {
	value := strings.TrimSpace(c.Transport.HTTP.Limit.Window)
	if value == "" {
		return time.Minute, nil
	}
	return time.ParseDuration(value)
}

type rawConfig struct {
	Database  rawDatabaseConfig  `yaml:"database"`
	Engine    rawEngineConfig    `yaml:"engine"`
	Transport rawTransportConfig `yaml:"transport"`
}

type rawDatabaseConfig struct {
	Path string `yaml:"path"`
}

type rawEngineConfig struct {
	Persistence rawEnginePersistenceConfig `yaml:"persistence"`
	Security    rawEngineSecurityConfig    `yaml:"security"`
}

type rawEnginePersistenceConfig struct {
	WindowSeconds int    `yaml:"windowSeconds"`
	WindowBytes   string `yaml:"windowBytes"`
	Threshold     string `yaml:"threshold"`
}

type rawEngineSecurityConfig struct {
	CompressAlgorithm string `yaml:"compressAlgorithm"`
	EncryptAlgorithm  string `yaml:"encryptAlgorithm"`
	EncryptKey        string `yaml:"encryptKey"`
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
	Limit         rawHTTPLimitConfig `yaml:"limit"`
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
	SQLGrantPath             string `yaml:"sqlGrantPath"`
	SQLRevokePath            string `yaml:"sqlRevokePath"`
	AdminPath                string `yaml:"adminPath"`
	ReportPath               string `yaml:"reportPath"`
	LimitEnabled             *bool  `yaml:"limitEnabled"`
	LimitRequests            int    `yaml:"limitRequests"`
	LimitWindow              string `yaml:"limitWindow"`
}

type rawHTTPLimitConfig struct {
	Enabled      *bool    `yaml:"enabled"`
	Requests     int      `yaml:"requests"`
	Window       string   `yaml:"window"`
	NoTokenPaths []string `yaml:"noTokenPaths"`
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
	SQLGrant             string `yaml:"sqlGrant"`
	SQLRevoke            string `yaml:"sqlRevoke"`
	Admin                string `yaml:"admin"`
	Report               string `yaml:"report"`
}

func applyRawConfig(config *Config, raw rawConfig) {
	if strings.TrimSpace(raw.Database.Path) != "" {
		config.Database.Path = strings.TrimSpace(raw.Database.Path)
	}
	if raw.Engine.Persistence.WindowSeconds > 0 {
		config.Engine.Persistence.WindowSeconds = raw.Engine.Persistence.WindowSeconds
	}
	if strings.TrimSpace(raw.Engine.Persistence.WindowBytes) != "" {
		config.Engine.Persistence.WindowBytes = strings.TrimSpace(raw.Engine.Persistence.WindowBytes)
	}
	if strings.TrimSpace(raw.Engine.Persistence.Threshold) != "" {
		config.Engine.Persistence.Threshold = strings.TrimSpace(raw.Engine.Persistence.Threshold)
	}
	if strings.TrimSpace(raw.Engine.Security.CompressAlgorithm) != "" {
		config.Engine.Security.CompressAlgorithm = strings.TrimSpace(raw.Engine.Security.CompressAlgorithm)
	}
	if strings.TrimSpace(raw.Engine.Security.EncryptAlgorithm) != "" {
		config.Engine.Security.EncryptAlgorithm = strings.TrimSpace(raw.Engine.Security.EncryptAlgorithm)
	}
	if strings.TrimSpace(raw.Engine.Security.EncryptKey) != "" {
		config.Engine.Security.EncryptKey = strings.TrimSpace(raw.Engine.Security.EncryptKey)
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
	if value := firstNonEmpty(raw.Transport.HTTP.Route.SQLGrant, raw.Transport.HTTP.SQLGrantPath); value != "" {
		config.Transport.HTTP.Route.SQLGrant = value
	}
	if value := firstNonEmpty(raw.Transport.HTTP.Route.SQLRevoke, raw.Transport.HTTP.SQLRevokePath); value != "" {
		config.Transport.HTTP.Route.SQLRevoke = value
	}
	if raw.Transport.HTTP.InitPassword != "" {
		config.Transport.HTTP.InitPassword = raw.Transport.HTTP.InitPassword
	}
	if strings.TrimSpace(raw.Transport.HTTP.TokenTTL) != "" {
		config.Transport.HTTP.TokenTTL = strings.TrimSpace(raw.Transport.HTTP.TokenTTL)
	}
	if raw.Transport.HTTP.Limit.Enabled != nil {
		config.Transport.HTTP.Limit.Enabled = *raw.Transport.HTTP.Limit.Enabled
	}
	if raw.Transport.HTTP.Limit.Requests > 0 {
		config.Transport.HTTP.Limit.Requests = raw.Transport.HTTP.Limit.Requests
	}
	if strings.TrimSpace(raw.Transport.HTTP.Limit.Window) != "" {
		config.Transport.HTTP.Limit.Window = strings.TrimSpace(raw.Transport.HTTP.Limit.Window)
	}
	if len(raw.Transport.HTTP.Limit.NoTokenPaths) > 0 {
		config.Transport.HTTP.Limit.NoTokenPaths = append([]string(nil), raw.Transport.HTTP.Limit.NoTokenPaths...)
	}
	if raw.Transport.HTTP.LimitEnabled != nil {
		config.Transport.HTTP.Limit.Enabled = *raw.Transport.HTTP.LimitEnabled
	}
	if raw.Transport.HTTP.LimitRequests > 0 {
		config.Transport.HTTP.Limit.Requests = raw.Transport.HTTP.LimitRequests
	}
	if strings.TrimSpace(raw.Transport.HTTP.LimitWindow) != "" {
		config.Transport.HTTP.Limit.Window = strings.TrimSpace(raw.Transport.HTTP.LimitWindow)
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

func (c Config) ParseEngineWindowBytes() (uint64, error) {
	c.ApplyDefaults()
	return parseBytes(c.Engine.Persistence.WindowBytes)
}

func (c Config) ParseEngineThresholdBytes() (uint64, error) {
	c.ApplyDefaults()
	return parseBytes(c.Engine.Persistence.Threshold)
}

func parseBytes(raw string) (uint64, error) {
	s := strings.TrimSpace(raw)
	s = strings.Trim(s, `"'`)
	if s == "" {
		return 0, fmt.Errorf("empty")
	}
	low := strings.ToLower(s)
	low = strings.ReplaceAll(low, "_", "")
	low = strings.ReplaceAll(low, " ", "")

	multiplier := uint64(1)
	switch {
	case strings.HasSuffix(low, "kb"):
		multiplier = 1024
		low = strings.TrimSuffix(low, "kb")
	case strings.HasSuffix(low, "k"):
		multiplier = 1024
		low = strings.TrimSuffix(low, "k")
	case strings.HasSuffix(low, "mb"):
		multiplier = 1024 * 1024
		low = strings.TrimSuffix(low, "mb")
	case strings.HasSuffix(low, "m"):
		multiplier = 1024 * 1024
		low = strings.TrimSuffix(low, "m")
	case strings.HasSuffix(low, "gb"):
		multiplier = 1024 * 1024 * 1024
		low = strings.TrimSuffix(low, "gb")
	case strings.HasSuffix(low, "g"):
		multiplier = 1024 * 1024 * 1024
		low = strings.TrimSuffix(low, "g")
	case strings.HasSuffix(low, "b"):
		low = strings.TrimSuffix(low, "b")
	}

	low = strings.TrimSpace(low)
	if low == "" {
		return 0, fmt.Errorf("missing number")
	}
	n, err := strconv.ParseUint(low, 10, 64)
	if err != nil {
		return 0, err
	}
	return n * multiplier, nil
}

func boolPtr(v bool) *bool { return &v }

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}
