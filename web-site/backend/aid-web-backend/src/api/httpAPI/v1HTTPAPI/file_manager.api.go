package v1HTTPAPI

import (
	"os"
	"path/filepath"

	"github.com/aid297/aid/filesystem/filesystemV4"
	"github.com/aid297/aid/validator/validatorV3"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/global"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/module/httpModule"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/module/httpModule/v1HTTPModule/request"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type FileManagerAPI struct{}

// Upload 上传单个文件
// * POST /api/v1/fileManger/upload
func (FileManagerAPI) Upload(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		httpModule.NewForbidden().SetMsg("获取上传文件失败").SetError(err).JSON(c)
		return
	}

	// 确保上传目录存在
	if err := os.MkdirAll(global.CONFIG.FileManager.Dir, 0755); err != nil {
		httpModule.NewInternalServerError().SetMsg("创建上传目录失败").SetError(err).JSON(c)
		return
	}

	savePath := filepath.Join(global.CONFIG.FileManager.Dir, file.Filename)

	// 保存文件到本地
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		httpModule.NewInternalServerError().SetMsg("保存文件失败").SetError(err).JSON(c)
		return
	}

	// 返回成功响应
	httpModule.NewOK().SetMsg("文件上传成功").SetData(gin.H{
		"filename":     file.Filename,
		"size":         file.Size,
		"saved_as":     file.Filename,
		"saved_path":   savePath,
		"content_type": file.Header.Get("Content-Type"),
	}).JSON(c)
}

// List 列出上传的文件
// * POST /api/v1/fileManger/list
func (FileManagerAPI) List(c *gin.Context) {
	var (
		tilte   = "获取文件列表"
		dir     = filesystemV4.NewDir(filesystemV4.Rel(global.CONFIG.FileManager.Dir))
		form    request.FileListRequest
		checker validatorV3.Checker
	)

	if form, checker = request.FileList.Bind(c); !checker.OK() {
		global.LOG.Error(tilte, zap.Any("表单验证", checker.Wrongs()))
		httpModule.NewForbidden().SetData(checker.Wrongs()).SetErrorf("表单验证失败：%w", checker.Wrong()).WithAccept(c)
		return
	}

	if form.Path != "" {
		dir = dir.Join(form.Path)
	}

	if !dir.GetExist() {
		dir.Create(filesystemV4.Mode(0777))
	}

	dir.LS()

	httpModule.NewOK().SetData(gin.H{"fullPath": dir.GetFullPath(), "dirs": dir.GetDirs(), "files": dir.GetFiles()}).WithAccept(c)
}
