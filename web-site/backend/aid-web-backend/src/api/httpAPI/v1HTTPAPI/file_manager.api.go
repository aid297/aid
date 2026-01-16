package v1HTTPAPI

import (
	"os"
	"path/filepath"

	"github.com/aid297/aid/filesystem/filesystemV3"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/global"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/module/httpModule"

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
		dir = filesystemV3.APP.Dir.Rel(filesystemV3.APP.DirAttr.Path.Set(global.CONFIG.FileManager.Dir))
	)

	if !dir.Exist {
		dir.Create(filesystemV3.DirMode(0777))
	}

	dir.LS()

	httpModule.NewOK().SetData(gin.H{"fullPath": dir.FullPath, "dirs": dir.Dirs, "files": dir.Files}).WithAccept(c)
}
