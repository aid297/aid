package httpAPI

import (
	"net/http"
	"time"

	"github.com/aid297/aid/web-site/backend/aid-web-backend-gin/global"
	"github.com/aid297/aid/web-site/backend/aid-web-backend-gin/module/httpModule"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HealthAPI struct{}

var healthRes = httpModule.HealthResponse{
	Time: httpModule.HealthTime{
		Now:    time.Now(),
		String: time.Now().Format(time.DateTime),
	},
	System: httpModule.HealthSystem{
		Debug:    global.CONFIG.System.Debug,
		Version:  global.CONFIG.System.Version,
		Daemon:   global.CONFIG.System.Daemon,
		Timezone: global.CONFIG.System.Timezone,
	},
	WebService: httpModule.HealthWebService{
		Cors: global.CONFIG.WebService.Cors,
	},
	VSCodeLaunch: httpModule.HealthVSCode{
		Version: "0.2.0",
		Configurations: []httpModule.HealthVSCodeConfiguration{
			{
				Name:    "Launch HA-BACKEND",
				Type:    "go",
				Request: "launch",
				Mode:    "auto",
				Program: "main.go",
				Env:     map[string]string{"HR_BACKEND_CONFIG": "config.json"},
				Args: []string{
					"-C=config.json",
					"-M=web-service",
					"-D=false",
				},
				BuildFlags: []string{"-tags=jsoniter"},
			},
		},
	},
}

// JSON 健康检查 JSON 格式
// * URL ANY /health/json
func (HealthAPI) JSON(c *gin.Context) {
	global.LOG.Info("获取支持的UUID版本", zap.Any("ANY /health/json", "生成UUID成功"))
	c.JSON(http.StatusOK, healthRes)
}

// JSON 健康检查 JSON 格式
// * URL ANY /health/yaml
func (HealthAPI) YAML(c *gin.Context) {
	global.LOG.Info("获取支持的UUID版本", zap.Any("ANY /health/yaml", "生成UUID成功"))
	c.YAML(http.StatusOK, healthRes)
}

// JSON 健康检查 JSON 格式
// * URL ANY /health/toml
func (HealthAPI) TOML(c *gin.Context) {
	global.LOG.Info("获取支持的UUID版本", zap.Any("ANY /health/toml", "生成UUID成功"))
	c.TOML(http.StatusOK, healthRes)
}
