package v1HTTPAPI

import (
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/aid297/aid/filesystem/filesystemV4"
	"github.com/aid297/aid/str"
	"github.com/aid297/aid/validator/validatorV3"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/global"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/module/httpModule"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/module/httpModule/v1HTTPModule/request"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/module/httpModule/v1HTTPModule/response"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type FileManagerAPI struct{}

// Upload 上传单个文件
// @Tags 文件管理
// @Summary 上传文件
// @Description 上传单个文件到指定路径
// @Produce application/json,application/xml,application/x-yaml,application/toml
// @Accept multipart/form-data
// @Router /fileManager/upload [post]
// @Param path query string true "上传路径"
// @Param file formData file true "上传文件内容"
// @Success 200 {object} httpModule.HTTPResponse{content=response.FileUploadResponse} "上传成功"
// @Failure 422 {object} httpModule.HTTPResponse "表单验证失败"
// @Failure 403 {object} httpModule.HTTPResponse "获取上传文件失败"
func (*FileManagerAPI) Upload(c *gin.Context) {
	var (
		title = "上传文件"
		err   error
		file  *multipart.FileHeader
		path  = c.Query("path")
	)

	// 获取上传的文件
	if file, err = c.FormFile("file"); err != nil {
		global.LOG.Error(title, zap.Errors("接收文件", []error{err}))
		httpModule.NewUnprocessableEntity(httpModule.Errorf("获取上传文件失败：%w", err)).WithAccept(c)
		return
	}

	// 确保上传目录存在
	if err = os.MkdirAll(filepath.Join(global.CONFIG.FileManager.Dir, path), 0755); err != nil {
		global.LOG.Error(title, zap.Errors("创建上传目录", []error{err}))
		httpModule.NewInternalServerError(httpModule.Errorf("创建上传目录失败：%w", err)).WithAccept(c)
		return
	}

	savePath := filepath.Join(global.CONFIG.FileManager.Dir, path, file.Filename)

	// 保存文件到本地
	if err = c.SaveUploadedFile(file, savePath); err != nil {
		global.LOG.Error(title, zap.Errors("保存文件到本地", []error{err}))
		httpModule.NewForbidden(httpModule.Errorf("保存文件到本地失败：%w", err)).WithAccept(c)
		return
	}

	global.LOG.Info(title, zap.Any("成功", ""))
	httpModule.NewCreated(
		httpModule.Msg("文件上传成功"),
		httpModule.Content(response.FileUploadResponse{
			FileName:    file.Filename,
			Size:        file.Size,
			ContentType: file.Header.Get("Content-Type"),
		}),
	).WithAccept(c)
}

// List 列出上传的文件
// @Tags 文件管理
// @Summary 获取文件列表
// @Description 按路径获取当前目录文件/目录列表
// @Produce application/json,application/xml,application/x-yaml,application/toml
// @Accept application/json
// @Param data body request.FileListRequest true "请求参数"
// @Router /fileManager/list [post]
// @Success 200 {object} httpModule.HTTPResponse{content=response.FileListResponse} "获取成功"
// @Failure 422 {object} httpModule.HTTPResponse "表单验证失败"
// @Failure 403 {object} httpModule.HTTPResponse "获取失败"
func (*FileManagerAPI) List(c *gin.Context) {
	type FilesystemerItem struct {
		Path string `json:"path"`
		Name string `json:"name"`
		Kind string `json:"kind"`
	}

	var (
		title         = "获取文件列表"
		err           error
		dir           filesystemV4.Filesystemer
		form          request.FileListRequest
		checker       validatorV3.Checker
		filesystemers []filesystemV4.Filesystemer
		rootPath      filesystemV4.Filesystemer
		currentPath   string
	)

	if form, checker = validatorV3.WithGin[request.FileListRequest](c); !checker.OK() {
		global.LOG.Error(title, zap.Any(global.ST_BIND_FORM, checker.Wrongs()))
		httpModule.NewUnprocessableEntity(httpModule.Content(checker.Wrongs()), httpModule.Errorf(global.FE_IVALIDED_FORM, checker.Wrong())).WithAccept(c)
		return
	}

	dir = filesystemV4.NewDir(filesystemV4.Rel(global.CONFIG.FileManager.Dir, form.Path))
	if form.Name == ".." && form.Path != "" {
		dir.Up()
	} else if form.Name != ".." && form.Name != "" {
		dir.Join(form.Name)
	}
	dir.LS()

	filesystemers = make([]filesystemV4.Filesystemer, 0, len(dir.GetDirs())+len(dir.GetFiles()))
	filesystemers = append(append(filesystemers, dir.GetDirs()...), dir.GetFiles()...)

	rootPath, err = filesystemV4.New(filesystemV4.Rel(global.CONFIG.FileManager.Dir))
	if err != nil {
		global.LOG.Error(title, zap.Errors("获取根路径错误", []error{err}))
		httpModule.NewInternalServerError(httpModule.Errorf("获取根路径错误：%w", err)).WithAccept(c)
		return
	}

	currentPath, _ = strings.CutPrefix(dir.GetFullPath(), rootPath.GetFullPath())

	global.LOG.Info(title, zap.Any("成功", filesystemers))
	httpModule.NewOK(httpModule.Content(response.FileListResponse{Filesystemers: filesystemers, CurrentPath: currentPath})).WithAccept(c)
}

// StoreFolder  创建文件夹
// @Tags        文件管理
// @Summary     创建文件夹
// @Produce     application/json
// @Accept      application/json
// @Param       data body request.FileStoreFolderRequest true "请求参数"
// @Router      /fileManager/storeFolder [post]
// @Success     200 {object} httpModule.HTTPResponse "创建成功"
// @Failure     422 {object} httpModule.HTTPResponse "表单验证失败"
// @Failure     403 {object} httpModule.HTTPResponse "创建文件夹失败"
func (*FileManagerAPI) StoreFolder(c *gin.Context) {
	var (
		title        = "创建文件夹"
		err          error
		filesystemer filesystemV4.Filesystemer
		form         request.FileStoreFolderRequest
		checker      validatorV3.Checker
	)

	if form, checker = validatorV3.WithGin[request.FileStoreFolderRequest](c); !checker.OK() {
		global.LOG.Error(title, zap.Any(global.ST_BIND_FORM, checker.Wrongs()))
		httpModule.NewUnprocessableEntity(httpModule.Content(checker.Wrongs()), httpModule.Errorf(global.FE_IVALIDED_FORM, checker.Wrong())).WithAccept(c)
		return
	}

	filesystemer = filesystemV4.NewDir(filesystemV4.Rel(global.CONFIG.FileManager.Dir, form.Path, form.Name))
	if err = filesystemer.Create(filesystemV4.Flag(0644)).GetError(); err != nil {
		global.LOG.Error(title, zap.Errors("创建文件夹失败", []error{err}))
		httpModule.NewForbidden(httpModule.Errorf("创建文件夹失败：%w", err)).WithAccept(c)
		return
	}

	global.LOG.Info(title, zap.String("成功", filesystemer.GetFullPath()))
	httpModule.NewOK(httpModule.Msg("创建成功")).WithAccept(c)
}

// Delete 删除文件或目录
// @Tags 文件管理
// @Summary 删除文件或目录
// @Produce application/json,application/xml,application/x-yaml,application/toml
// @Accept application/json
// @Param data body request.FileDestroyRequest true "请求参数"
// @Router /api/v1/fileManger/destroy [post]
// @Success 200 {object} httpModule.HTTPResponse "删除成功"
// @Failure 422 {object} httpModule.HTTPResponse "表单验证失败"
// @Failure 403 {object} httpModule.HTTPResponse "删除失败"
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
		httpModule.NewUnprocessableEntity(httpModule.Content(checker.Wrongs()), httpModule.Errorf(global.FE_IVALIDED_FORM, checker.Wrong())).WithAccept(c)
		return
	}

	if filesystemer, err = filesystemV4.New(filesystemV4.Rel(global.CONFIG.FileManager.Dir, form.Path, form.Name)); err != nil {
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
// @Tags 文件管理
// @Summary 下载文件
// @Produce application/octet-stream
// @Accept application/json
// @Param path query string true "文件路径"
// @Param name query string true "文件名"
// @Router /fileManager/download [get]
// @Success 200 {file} file "下载成功"
// @Failure 422 {object} httpModule.HTTPResponse "表单验证失败"
// @Failure 403 {object} httpModule.HTTPResponse "下载失败"
func (*FileManagerAPI) Download(c *gin.Context) {
	var (
		dir  filesystemV4.Filesystemer
		path = c.Query("path")
		name = c.Query("name")
	)

	if dir = filesystemV4.NewDir(filesystemV4.Rel(global.CONFIG.FileManager.Dir, path, name)); !dir.GetExist() {
		c.String(404, str.APP.HTML.New(
			str.HtmlH(1, str.HtmlNormal("错误：文件不存在")),
			str.HtmlH(2, str.HtmlNormal(dir.GetFullPath())),
		).End())
		return
	}

	// 设置文件名
	c.Header("Content-Disposition", "attachment; filename="+name)
	c.File(dir.GetFullPath())
}

// Zip 压缩文件或目录
// @Tags 文件管理
// @Summary 压缩文件或目录
// @Produce application/json,application/xml,application/x-yaml,application/toml
// @Accept application/json
// @Param data body request.FileZipRequest true "请求参数"
// @Router /fileManager/zip [post]
// @Success 200 {object} httpModule.HTTPResponse{content=response.FileZipResponse} "压缩成功"
// @Failure 422 {object} httpModule.HTTPResponse "表单验证失败"
// @Failure 403 {object} httpModule.HTTPResponse "压缩失败"
// @Failure 404 {object} httpModule.HTTPResponse "获取路径错误"
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
		httpModule.NewUnprocessableEntity(httpModule.Content(checker.Wrongs()), httpModule.Errorf(global.FE_IVALIDED_FORM, checker.Wrong())).WithAccept(c)
		return
	}

	if filesystemer, err = filesystemV4.New(filesystemV4.Rel(global.CONFIG.FileManager.Dir, form.Path, form.Name)); err != nil {
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
	httpModule.NewOK(httpModule.Msg("压缩成功"), httpModule.Content(response.FileZipResponse{Name: zipped.GetName()})).WithAccept(c)
}
