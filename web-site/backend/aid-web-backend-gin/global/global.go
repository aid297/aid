package global

import (
	"github.com/aid297/aid/web-site/backend/aid-web-backend-gin/config"

	"go.uber.org/zap"
)

var (
	CONFIG config.Config
	LOG    *zap.Logger
)
