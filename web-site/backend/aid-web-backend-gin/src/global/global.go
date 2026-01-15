package global

import (
	"go.uber.org/zap"

	`github.com/aid297/aid/web-site/backend/aid-web-backend-gin/src/config`
)

var (
	CONFIG config.Config
	LOG    *zap.Logger
)
