package v1HTTPAPI

import (
	"github.com/aid297/aid/validator/validatorV3"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/global"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/module/httpModule"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/module/httpModule/v1HTTPModule/request"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/module/httpModule/v1HTTPModule/response"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/service/httpService/v1HTTPService"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UUIDAPI struct{}

// Generate 批量生成uuid
// * URL POST /api/v1/uuid/generate
func (*UUIDAPI) Generate(c *gin.Context) {
	var (
		title   = "批量生成uuid"
		err     error
		form    request.UUIDGenerateRequest
		checker validatorV3.Checker
		uuids   []response.UUIDResponse
	)

	if form, checker = validatorV3.WithGin[request.UUIDGenerateRequest](c); !checker.OK() {
		global.LOG.Error(title, zap.Any(global.ST_BIND_FORM, checker.Wrongs()))
		httpModule.NewForbidden(httpModule.Content(checker.Wrongs()), httpModule.Errorf(global.FE_IVALIDED_FORM, checker.Wrong())).WithAccept(c)
		return
	}

	uuids = make([]response.UUIDResponse, form.Number)

	for idx := range form.Number {
		if uuids[idx], err = v1HTTPService.UUID.GenerateOne(&form); err != nil {
			global.LOG.Error(title, zap.Any("生成UUID失败", err.Error()))
			httpModule.NewForbidden(httpModule.Errorf("生成UUID失败：%w", err)).WithAccept(c)
			return
		}
	}

	global.LOG.Info(title, zap.Any("POST /api/v1/uuid/generate", "生成UUID成功"))
	httpModule.NewOK(httpModule.Content(response.UUIDGenerateResponse{UUIDs: uuids})).WithAccept(c)
}

// Versions 获取支持的UUID版本
// * URL POST /api/v1/uuid/versions
func (*UUIDAPI) Versions(c *gin.Context) {
	global.LOG.Info("获取支持的UUID版本", zap.Any("POST /api/v1/uuid/versions", "生成UUID成功"))
	httpModule.NewOK(httpModule.Content(response.UUIDVersionsResponse{
		Versions: map[string]string{
			string(request.UUIDVersionV1): string(request.UUIDVersionV1),
			string(request.UUIDVersionV4): string(request.UUIDVersionV4),
			string(request.UUIDVersionV6): string(request.UUIDVersionV6),
			string(request.UUIDVersionV7): string(request.UUIDVersionV7),
		},
	})).WithAccept(c)
}
