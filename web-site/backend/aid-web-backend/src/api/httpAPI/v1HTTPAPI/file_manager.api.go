package v1HTTPAPI

import (
	"mime/multipart"
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
func (*FileManagerAPI) Upload(c *gin.Context) {
	var (
		title = "上传文件"
		err   error
		file  *multipart.FileHeader
	)

	// 获取上传的文件
	if file, err = c.FormFile("file"); err != nil {
		global.LOG.Error(title, zap.Errors("接收文件", []error{err}))
		httpModule.NewForbidden(httpModule.Errorf("获取上传文件失败：%w", err)).JSON(c)
		return
	}

	// 确保上传目录存在
	if err := os.MkdirAll(global.CONFIG.FileManager.Dir, 0755); err != nil {
		global.LOG.Error(title, zap.Errors("创建上传目录", []error{err}))
		httpModule.NewInternalServerError(httpModule.Errorf("创建上传目录失败：%w", err)).JSON(c)
		return
	}

	savePath := filepath.Join(global.CONFIG.FileManager.Dir, file.Filename)

	// 保存文件到本地
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		global.LOG.Error(title, zap.Errors("保存文件到本地", []error{err}))
		httpModule.NewForbidden(httpModule.Errorf("保存文件到本地失败：%w", err)).JSON(c)
		return
	}

	global.LOG.Info(title, zap.Any("成功", ""))
	httpModule.NewOK(
		httpModule.Msg("文件上传成功"),
		httpModule.Content(gin.H{
			"filename":     file.Filename,
			"size":         file.Size,
			"saved_as":     file.Filename,
			"saved_path":   savePath,
			"content_type": file.Header.Get("Content-Type"),
		}),
	).WithAccept(c)
}

// List 列出上传的文件
// * POST /api/v1/fileManger/list
func (*FileManagerAPI) List(c *gin.Context) {
	var (
		title   = "获取文件列表"
		dir     filesystemV4.Filesystemer
		form    request.FileListRequest
		checker validatorV3.Checker
	)

	if form, checker = validatorV3.WithGin[request.FileListRequest](c); !checker.OK() {
		global.LOG.Error(title, zap.Any(global.ST_BIND_FORM, checker.Wrongs()))
		httpModule.NewForbidden(httpModule.Content(checker.Wrongs()), httpModule.Errorf(global.FE_IVALIDED_FORM, checker.Wrong())).WithAccept(c)
		return
	}

	dir = filesystemV4.NewDir(filesystemV4.Rel(global.CONFIG.FileManager.Dir, form.Path))
	if !dir.GetExist() {
		httpModule.NewNotFound(httpModule.Msg("文件管理系统不存在")).WithAccept(c)
		return
	}

	dir.LS()

	ret := gin.H{"fullPath": dir.GetFullPath(), "dirs": dir.GetDirs(), "files": dir.GetFiles()}
	global.LOG.Info(title, zap.Any("成功", ret))
	httpModule.NewOK(httpModule.Content(ret)).WithAccept(c)
}

// Delete 删除文件或目录
// * POST /api/v1/fileManger/destroy
func (*FileManagerAPI) Destroy(c *gin.Context) {
	var (
		title        = "删除文件或目录"
		err          error
		filesystemer filesystemV4.Filesystemer
		form         request.FileDestroyRequest
		checker      validatorV3.Checker
	)

	if form, checker = validatorV3.WithGin[request.FileDestroyRequest](c); !checker.OK() {
		global.LOG.Error(title, zap.Any(global.ST_BIND_FORM, checker.Wrongs()))
		httpModule.NewForbidden(httpModule.Content(checker.Wrongs()), httpModule.Errorf(global.FE_IVALIDED_FORM, checker.Wrong())).WithAccept(c)
		return
	}

	if filesystemer, err = filesystemV4.New(filesystemV4.Rel(global.CONFIG.FileManager.Dir, form.Path)); err != nil {
		global.LOG.Error(title, zap.Errors("获取路径错误", []error{err}))
		httpModule.NewNotFound(httpModule.Errorf("获取路径错误：%w", err)).WithAccept(c)
		return
	}

	if err = filesystemer.RemoveAll().GetError(); err != nil {
		global.LOG.Error(title, zap.Errors("删除文件或目录失败", []error{err}))
		httpModule.NewForbidden(httpModule.Errorf("删除文件或目录失败：%w", err)).WithAccept(c)
		return
	}

	global.LOG.Info(title, zap.String("成功", form.Path))
	httpModule.NewOK(httpModule.Msg("删除成功")).WithAccept(c)
}

// Download 下载文件
// * POST /api/v1/fileManger/download
func (*FileManagerAPI) Download(c *gin.Context) {
	var (
		title   = "下载文件"
		dir     filesystemV4.Filesystemer
		form    request.FileDownloadRequest
		checker validatorV3.Checker
	)

	if form, checker = validatorV3.WithGin[request.FileDownloadRequest](c); !checker.OK() {
		global.LOG.Error(title, zap.Any(global.ST_BIND_FORM, checker.Wrongs()))
		httpModule.NewForbidden(httpModule.Content(checker.Wrongs()), httpModule.Errorf(global.FE_IVALIDED_FORM, checker.Wrong())).WithAccept(c)
		return
	}

	dir = filesystemV4.NewDir(filesystemV4.Rel(global.CONFIG.FileManager.Dir, form.Path))
	if !dir.GetExist() {
		httpModule.NewNotFound(httpModule.Msg("文件不存在")).WithAccept(c)
		return
	}

	c.File(dir.GetFullPath())
}

// Zip 压缩文件或目录
// * POST /api/v1/fileManger/zip
func (*FileManagerAPI) Zip(c *gin.Context) {
	var (
		title                = "压缩文件或目录"
		err                  error
		filesystemer, zipped filesystemV4.Filesystemer
		form                 request.FileZipRequest
		checker              validatorV3.Checker
	)

	if form, checker = validatorV3.WithGin[request.FileZipRequest](c); !checker.OK() {
		global.LOG.Error(title, zap.Any(global.ST_BIND_FORM, checker.Wrongs()))
		httpModule.NewForbidden(httpModule.Content(checker.Wrongs()), httpModule.Errorf(global.FE_IVALIDED_FORM, checker.Wrong())).WithAccept(c)
		return
	}

	if filesystemer, err = filesystemV4.New(filesystemV4.Rel(global.CONFIG.FileManager.Dir, form.Path)); err != nil {
		global.LOG.Error(title, zap.Errors("获取路径错误", []error{err}))
		httpModule.NewNotFound(httpModule.Errorf("获取路径错误：%w", err)).WithAccept(c)
		return
	}

	zipped = filesystemer.Zip()
	if err = zipped.GetError(); err != nil {
		global.LOG.Error(title, zap.Errors("压缩文件或目录失败", []error{err}))
		httpModule.NewForbidden(httpModule.Errorf("压缩文件或目录失败：%w", err)).WithAccept(c)
		return
	}

	global.LOG.Info(title, zap.String("成功", zipped.GetFullPath()))
	httpModule.NewOK(httpModule.Msg("压缩成功"), httpModule.Content(gin.H{"name": zipped.GetName()})).WithAccept(c)
}
