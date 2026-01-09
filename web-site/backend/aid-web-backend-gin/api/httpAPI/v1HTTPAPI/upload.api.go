package v1HTTPAPI

import (
	"os"
	"path/filepath"

	"github.com/aid297/aid/web-site/backend/aid-web-backend-gin/global"
	"github.com/aid297/aid/web-site/backend/aid-web-backend-gin/module/httpModule"

	"github.com/gin-gonic/gin"
)

type UploadAPI struct{}

// Single 上传单个文件
// * Single /api/v1/upload/single
func (UploadAPI) Single(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		httpModule.NewForbidden().SetMsg("获取上传文件失败").SetError(err).JSON(c)
		return
	}

	// 确保上传目录存在
	if err := os.MkdirAll(global.CONFIG.UploadDir, 0755); err != nil {
		httpModule.NewInternalServerError().SetMsg("创建上传目录失败").SetError(err).JSON(c)
		return
	}

	savePath := filepath.Join(global.CONFIG.UploadDir, file.Filename)

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

// Multiple 上传多个文件
// * POST /api/v1/upload/multiple
func (UploadAPI) Multiple(c *gin.Context) {
	// 获取多个上传文件
	form, err := c.MultipartForm()
	if err != nil {
		httpModule.NewForbidden().SetMsg("获取上传文件失败").SetError(err).JSON(c)
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		httpModule.NewForbidden().SetMsg("没有上传任何文件").JSON(c)
		return
	}

	// 确保上传目录存在
	if err := os.MkdirAll(global.CONFIG.UploadDir, 0755); err != nil {
		httpModule.NewInternalServerError().SetMsg("创建上传目录失败").SetError(err).JSON(c)
		return
	}

	// 保存所有文件
	var uploadedFiles []gin.H
	for _, file := range files {
		// 生成唯一文件名
		savePath := filepath.Join(global.CONFIG.UploadDir, file.Filename)

		// 保存文件
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			httpModule.NewInternalServerError().SetMsg("保存文件失败: " + file.Filename).SetError(err).JSON(c)
			return
		}

		uploadedFiles = append(uploadedFiles, gin.H{
			"filename":     file.Filename,
			"size":         file.Size,
			"saved_as":     file.Filename,
			"saved_path":   savePath,
			"content_type": file.Header.Get("Content-Type"),
		})
	}

	// 返回成功响应
	httpModule.NewOK().SetMsg("文件上传成功").SetData(gin.H{
		"count": len(uploadedFiles),
		"files": uploadedFiles,
	}).JSON(c)
}
