package main

import (
	"flag"

	"github.com/aid297/aid/debugLogger"
	configpkg "github.com/aid297/aid/simpleDB/config"
	"github.com/aid297/aid/simpleDB/transport"
	"github.com/gin-gonic/gin"
)

func main() {
	defaults := configpkg.Default()
	var (
		addr      = flag.String("addr", defaults.Transport.HTTP.Address, "HTTP listen address")
		database  = flag.String("db", defaults.Database.Path, "simpleDB database path or name")
		loginPath = flag.String("login-path", defaults.Transport.HTTP.Route.Login, "login route path")
	)
	flag.Parse()

	server := transport.New.HTTP(*database, transport.WithLoginPath(*loginPath))
	registerDemoRoutes(server)
	printStartupInfo(*addr, *database, *loginPath)

	if err := server.Run(*addr); err != nil {
		debugLogger.Fatal("transport http server stopped: %v", err)
	}
}

func registerDemoRoutes(server *transport.HTTPServer) {
	server.Engine().GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]any{
			"success": true,
			"service": "simpleDB transport",
			"status":  "ok",
		})
	})

	server.Engine().GET("/me", server.AuthMiddleware(), func(ctx *gin.Context) {
		claims, ok := transport.UserFromContext(ctx)
		if !ok {
			ctx.JSON(500, map[string]any{"success": false, "error": "missing auth context"})
			return
		}
		ctx.JSON(200, map[string]any{"success": true, "user": claims})
	})

	server.Engine().GET("/admin", server.AuthMiddleware(), server.RequireRoles("super_admin"), func(ctx *gin.Context) {
		ctx.JSON(200, map[string]any{"success": true, "message": "role check passed"})
	})

	server.Engine().GET("/reports", server.AuthMiddleware(), server.RequirePermissions("report.read"), func(ctx *gin.Context) {
		ctx.JSON(200, map[string]any{"success": true, "message": "permission check passed"})
	})
}

func printStartupInfo(addr, database, loginPath string) {
	debugLogger.Print("simpleDB transport HTTP demo is running")
	debugLogger.Print("- addr: %s", addr)
	debugLogger.Print("- database: %s\n", database)
	debugLogger.Print("- login: POST %s\n", loginPath)
	debugLogger.Print("- health: GET /health\n")
	debugLogger.Print("- me: GET /me (requires Bearer token)\n")
	debugLogger.Print("- admin: GET /admin (requires role super_admin)\n")
	debugLogger.Print("- reports: GET /reports (requires permission report.read)\n")
	debugLogger.Print("- refresh: POST %s (requires Bearer token)\n", transport.DefaultRefreshPath)
	debugLogger.Print("- logout: POST %s (requires Bearer token)\n", transport.DefaultLogoutPath)
	debugLogger.Print("- default admin: sdb / simpleDB\n\n")
	debugLogger.Print("curl example:\n")
	debugLogger.Print("curl -X POST http://%s%s \\\n", normalizeAddr(addr), loginPath)
	debugLogger.Print("  -H 'Content-Type: application/json' \\\n")
	debugLogger.Print("  -d '{\"username\":\"sdb\",\"password\":\"simpleDB\"}'\n")
	debugLogger.Print("\nThen call /me with the returned token:\n")
	debugLogger.Print("curl http://%s/me -H 'Authorization: Bearer <accessToken>'\n", normalizeAddr(addr))
	debugLogger.Print("curl -X POST http://%s%s -H 'Authorization: Bearer <accessToken>'\n", normalizeAddr(addr), transport.DefaultRefreshPath)
	debugLogger.Print("curl -X POST http://%s%s -H 'Authorization: Bearer <accessToken>'\n", normalizeAddr(addr), transport.DefaultLogoutPath)
}

func normalizeAddr(addr string) string {
	if addr == "" {
		return "127.0.0.1:18080"
	}
	if addr[0] == ':' {
		return "127.0.0.1" + addr
	}
	return addr
}
