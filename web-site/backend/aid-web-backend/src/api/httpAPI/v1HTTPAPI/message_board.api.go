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

// MesasgeBoardAPI API：留言板
type MessageBoardAPI struct{}

// List 留言板API：列表
// @Tags 留言板
// @Summary 获取留言板信息列表
// @Produce application/json,application/xml,application/x-yaml,application/toml
// @Accept application/json
// @Router /messageBoard/list [post]
// @Success 200 {object} httpModule.HTTPResponse{content=response.MessageBoardListResponse} "获取留言板信息列表成功"
// @Failure 422 {object} httpModule.HTTPResponse "表单验证失败"
// @Failure 403 {object} httpModule.HTTPResponse "获取留言板信息列表失败"
// @Failure 404 {object} httpModule.HTTPResponse "获取路径错误"
func (*MessageBoardAPI) List(c *gin.Context) {
	var (
		title = "获取留言板信息列表"
		err   error
		res   response.MessageBoardListResponse
	)

	if res.MessageBoards, err = v1HTTPService.New.MessageBoard().List(); err != nil {
		global.LOG.Error(title, zap.Any("获取留言板信息列表失败", err.Error()))
		httpModule.NewForbidden(httpModule.Errorf("获取留言板信息列表失败：%w", err)).WithAccept(c)
		return
	}

	httpModule.NewOK(httpModule.Content(res)).WithAccept(c)
}

// Store 留言板API：保存
// @Tags 留言板
// @Summary 保存留言板信息
// @Produce application/json,application/xml,application/x-yaml,application/toml
// @Accept application/json
// @Param data body request.MessageBoardStoreRequest true "请求参数"
// @Router /messageBoard/store [post]
// @Success 200 {object} httpModule.HTTPResponse "保存留言板信息成功"
// @Failure 422 {object} httpModule.HTTPResponse "表单验证失败"
// @Failure 403 {object} httpModule.HTTPResponse "保存留言板信息失败"
// @Failure 404 {object} httpModule.HTTPResponse "获取路径错误"
func (*MessageBoardAPI) Store(c *gin.Context) {
	var (
		title   = "保存留言板信息"
		err     error
		form    request.MessageBoardStoreRequest
		checker validatorV3.Checker
	)

	if form, checker = (&request.MessageBoardStoreRequest{}).Bind(c); !checker.OK() {
		global.LOG.Error(title, zap.Any(global.ST_BIND_FORM, checker.Wrongs()))
		httpModule.NewForbidden(httpModule.Content(checker.Wrongs()), httpModule.Errorf(global.FE_IVALIDED_FORM, checker.Wrong())).WithAccept(c)
		return
	}

	if err = v1HTTPService.New.MessageBoard().Store(&form); err != nil {
		global.LOG.Error(title, zap.Any("保存留言板信息失败", err.Error()))
		httpModule.NewForbidden(httpModule.Errorf("保存留言板信息失败：%w", err)).WithAccept(c)
		return
	}

	httpModule.NewOK().WithAccept(c)
}

// Destroy 留言板API：删除
// @Tags 留言板
// @Summary 删除留言板信息
// @Produce application/json,application/xml,application/x-yaml,application/toml
// @Accept application/json
// @Param data body request.MessageBoardDestroyRequest true "请求参数"
// @Router /messageBoard/destroy [post]
// @Success 200 {object} httpModule.HTTPResponse "删除留言板信息成功"
// @Failure 422 {object} httpModule.HTTPResponse "表单验证失败"
// @Failure 403 {object} httpModule.HTTPResponse "删除留言板信息失败"
// @Failure 404 {object} httpModule.HTTPResponse "获取路径错误"
func (*MessageBoardAPI) Destroy(c *gin.Context) {
	var (
		title   = "删除留言板信息"
		err     error
		form    request.MessageBoardDestroyRequest
		checker validatorV3.Checker
	)

	if form, checker = (&request.MessageBoardDestroyRequest{}).Bind(c); !checker.OK() {
		global.LOG.Error(title, zap.Any(global.ST_BIND_FORM, checker.Wrongs()))
		httpModule.NewForbidden(httpModule.Content(checker.Wrongs()), httpModule.Errorf(global.FE_IVALIDED_FORM, checker.Wrong())).WithAccept(c)
		return
	}

	if err = v1HTTPService.New.MessageBoard().Destroy(&form); err != nil {
		global.LOG.Error(title, zap.Any("删除留言板信息失败", err.Error()))
		httpModule.NewForbidden(httpModule.Errorf("删除留言板信息失败：%w", err)).WithAccept(c)
		return
	}

	httpModule.NewDeleted(httpModule.Msg("删除成功")).WithAccept(c)
}
