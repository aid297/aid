package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	configpkg "github.com/aid297/aid/simpleDB/config"
	"github.com/aid297/aid/simpleDB/transport"
	"github.com/gin-gonic/gin"
)

func runServe(args []string, stdout, stderr *os.File) int {
	fs := newFlagSet("serve")
	configPath := fs.String("config", defaultConfigPath(), "config file path")
	if err := fs.Parse(args); err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		return 2
	}

	resolvedConfigPath := resolveConfigPath(*configPath)

	config, err := configpkg.Load(resolvedConfigPath)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "load config failed: %v\n", err)
		return 1
	}
	if err = config.Validate(); err != nil {
		_, _ = fmt.Fprintf(stderr, "invalid config: %v\n", err)
		return 1
	}

	if err = ensureInitPassword(resolvedConfigPath, &config, stdout); err != nil {
		_, _ = fmt.Fprintf(stderr, "ensure init password failed: %v\n", err)
		return 1
	}

	server, err := newHTTPServerFromConfig(config, resolvedConfigPath)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "build HTTP server failed: %v\n", err)
		return 1
	}

	registerServiceRoutes(server, config)
	printServeSummary(stdout, resolvedConfigPath, config)

	if err = server.Run(config.Transport.HTTP.Address); err != nil {
		_, _ = fmt.Fprintf(stderr, "HTTP service stopped: %v\n", err)
		return 1
	}
	return 0
}

func newHTTPServerFromConfig(config configpkg.Config, configPath string) (*transport.HTTPServer, error) {
	tokenTTL, err := config.ParseTokenTTL()
	if err != nil {
		return nil, err
	}
	limitWindow, err := config.ParseLimitWindow()
	if err != nil {
		return nil, err
	}
	resolvedConfigPath := resolveConfigPath(configPath)
	gin.SetMode(config.Transport.HTTP.GinMode)
	server := transport.New.HTTP(
		config.Database.Path,
		transport.WithLoginPath(config.Transport.HTTP.Route.Login),
		transport.WithRegisterPath(config.Transport.HTTP.Route.Register),
		transport.WithRefreshPath(config.Transport.HTTP.Route.Refresh),
		transport.WithLogoutPath(config.Transport.HTTP.Route.Logout),
		transport.WithActivatePath(config.Transport.HTTP.Route.Activate),
		transport.WithDeactivatePath(config.Transport.HTTP.Route.Deactivate),
		transport.WithAssignRolePath(config.Transport.HTTP.Route.AssignRole),
		transport.WithAssignRolePermissionPath(config.Transport.HTTP.Route.AssignRolePermission),
		transport.WithInitSDBPasswordPath(config.Transport.HTTP.Route.InitSDBPassword),
		transport.WithSQLExecutePath(config.Transport.HTTP.Route.SQLExecute),
		transport.WithSQLGrantPath(config.Transport.HTTP.Route.SQLGrant),
		transport.WithSQLRevokePath(config.Transport.HTTP.Route.SQLRevoke),
		transport.WithSQLAllowedOps(config.Transport.HTTP.SQLAllowedOps),
		transport.WithTokenRateLimit(
			config.Transport.HTTP.Limit.Enabled,
			config.Transport.HTTP.Limit.Requests,
			limitWindow,
			config.Transport.HTTP.Limit.NoTokenPaths,
		),
		transport.WithInitPassword(config.Transport.HTTP.InitPassword),
		transport.WithInitPasswordRotator(func() (string, error) {
			return rotateInitPasswordInConfig(resolvedConfigPath)
		}),
		transport.WithTokenTTL(tokenTTL),
		transport.WithTokenSecret(config.Transport.HTTP.TokenSecret),
	)
	return server, nil
}

func registerServiceRoutes(server *transport.HTTPServer, config configpkg.Config) {
	server.Engine().GET(config.Transport.HTTP.Route.Health, func(ctx *gin.Context) {
		ctx.JSON(200, map[string]any{"success": true, "service": "simpleDB", "status": "ok"})
	})

	server.Engine().GET(config.Transport.HTTP.Route.Profile, server.AuthMiddleware(), func(ctx *gin.Context) {
		claims, ok := transport.UserFromContext(ctx)
		if !ok {
			ctx.JSON(500, map[string]any{"success": false, "error": "missing auth context"})
			return
		}
		ctx.JSON(200, map[string]any{"success": true, "user": claims})
	})

	if config.Transport.HTTP.EnableAdmin {
		server.Engine().GET(config.Transport.HTTP.Route.Admin, server.AuthMiddleware(), server.RequireRoles("super_admin"), func(ctx *gin.Context) {
			ctx.JSON(200, map[string]any{"success": true, "message": "role check passed"})
		})
	}

	if config.Transport.HTTP.EnableReport {
		server.Engine().GET(config.Transport.HTTP.Route.Report, server.AuthMiddleware(), server.RequirePermissions("report.read"), func(ctx *gin.Context) {
			ctx.JSON(200, map[string]any{"success": true, "message": "permission check passed"})
		})
	}
}

func printServeSummary(stdout *os.File, configPath string, config configpkg.Config) {
	_, _ = fmt.Fprintln(stdout, "simpleDB service starting")
	_, _ = fmt.Fprintf(stdout, "- config: %s\n", configPath)
	_, _ = fmt.Fprintf(stdout, "- database: %s\n", config.Database.Path)
	_, _ = fmt.Fprintf(stdout, "- http: %s\n", config.Transport.HTTP.Address)
	_, _ = fmt.Fprintf(stdout, "- login: POST %s\n", config.Transport.HTTP.Route.Login)
	_, _ = fmt.Fprintf(stdout, "- register: POST %s\n", config.Transport.HTTP.Route.Register)
	_, _ = fmt.Fprintf(stdout, "- refresh: POST %s\n", config.Transport.HTTP.Route.Refresh)
	_, _ = fmt.Fprintf(stdout, "- logout: POST %s\n", config.Transport.HTTP.Route.Logout)
	_, _ = fmt.Fprintf(stdout, "- activate: POST %s (requires super_admin)\n", config.Transport.HTTP.Route.Activate)
	_, _ = fmt.Fprintf(stdout, "- deactivate: POST %s (requires super_admin)\n", config.Transport.HTTP.Route.Deactivate)
	_, _ = fmt.Fprintf(stdout, "- assign roles: POST %s (requires super_admin)\n", config.Transport.HTTP.Route.AssignRole)
	_, _ = fmt.Fprintf(stdout, "- assign role permissions: POST %s (requires super_admin)\n", config.Transport.HTTP.Route.AssignRolePermission)
	_, _ = fmt.Fprintf(stdout, "- init sdb password: POST %s (requires init password)\n", config.Transport.HTTP.Route.InitSDBPassword)
	_, _ = fmt.Fprintf(stdout, "- sql execute: POST %s (requires auth)\n", config.Transport.HTTP.Route.SQLExecute)
	_, _ = fmt.Fprintf(stdout, "- sql grant: POST %s (requires auth)\n", config.Transport.HTTP.Route.SQLGrant)
	_, _ = fmt.Fprintf(stdout, "- sql revoke: POST %s (requires auth)\n", config.Transport.HTTP.Route.SQLRevoke)
	if config.Transport.HTTP.Limit.Enabled {
		_, _ = fmt.Fprintf(stdout, "- limit: enabled (%d requests / %s)\n", config.Transport.HTTP.Limit.Requests, config.Transport.HTTP.Limit.Window)
	} else {
		_, _ = fmt.Fprintln(stdout, "- limit: disabled")
	}
	_, _ = fmt.Fprintf(stdout, "- health: GET %s\n", config.Transport.HTTP.Route.Health)
	_, _ = fmt.Fprintf(stdout, "- profile: GET %s\n", config.Transport.HTTP.Route.Profile)
	if config.Transport.HTTP.EnableAdmin {
		_, _ = fmt.Fprintf(stdout, "- admin: GET %s\n", config.Transport.HTTP.Route.Admin)
	}
	if config.Transport.HTTP.EnableReport {
		_, _ = fmt.Fprintf(stdout, "- reports: GET %s\n", config.Transport.HTTP.Route.Report)
	}
	_, _ = fmt.Fprintln(stdout, "- default admin: sdb / simpleDB")
}

func ensureInitPassword(configPath string, config *configpkg.Config, stdout *os.File) error {
	if config == nil {
		return fmt.Errorf("config is nil")
	}
	if strings.TrimSpace(config.Transport.HTTP.InitPassword) != "" {
		return nil
	}

	generated, err := generateRandomString(24)
	if err != nil {
		return err
	}
	config.Transport.HTTP.InitPassword = generated

	path := resolveConfigPath(configPath)
	if err = configpkg.Save(path, *config); err != nil {
		return err
	}

	if stdout != nil {
		_, _ = fmt.Fprintf(stdout, "- generated init password and saved to config: %s\n", path)
	}
	return nil
}

func resolveConfigPath(path string) string {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		trimmed = defaultConfigPath()
	}

	if filepath.IsAbs(trimmed) {
		return trimmed
	}

	if _, err := os.Stat(trimmed); err == nil {
		return trimmed
	}

	const prefix = "simpleDB/"
	if after, ok := strings.CutPrefix(trimmed, prefix); ok {
		fallback := after
		if _, err := os.Stat(fallback); err == nil {
			return fallback
		}
	}

	return trimmed
}

func generateRandomString(length int) (string, error) {
	if length <= 0 {
		length = 24
	}
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	text := base64.RawURLEncoding.EncodeToString(buf)
	if len(text) >= length {
		return text[:length], nil
	}
	return text, nil
}

func rotateInitPasswordInConfig(configPath string) (string, error) {
	path := resolveConfigPath(configPath)
	config, err := configpkg.Load(path)
	if err != nil {
		return "", err
	}

	generated, err := generateRandomString(24)
	if err != nil {
		return "", err
	}
	config.Transport.HTTP.InitPassword = generated

	if err = configpkg.Save(path, config); err != nil {
		return "", err
	}

	return generated, nil
}
