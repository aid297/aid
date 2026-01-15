package v1HTTPAPI

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	`github.com/aid297/aid/web-site/backend/aid-web-backend-gin/src/global`
	`github.com/aid297/aid/web-site/backend/aid-web-backend-gin/src/module/httpModule`

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/aid297/aid/filesystem/filesystemV2"
)

type RezipAPI struct{}

func (RezipAPI) Upload(c *gin.Context) {
	var (
		title     = "重新压缩"
		file      *multipart.FileHeader
		err       error
		srcFile   multipart.File
		buf       []byte
		zipReader *zip.Reader
		rc        io.ReadCloser
	)

	if file, err = c.FormFile("f"); err != nil {
		global.LOG.Error(title, zap.Errors("上传文件", []error{err}))
		httpModule.NewForbidden().SetMsg("上传文件失败").WithAccept(c)
		return
	}

	// 打开上传文件
	srcFile, err = file.Open()
	if err != nil {
		global.LOG.Error(title, zap.Errors("打开文件", []error{err}))
		httpModule.NewForbidden().SetErrorf("打开文件失败：%w", err).WithAccept(c)
		return
	}
	defer func() { _ = srcFile.Close() }()

	// 读到内存
	buf = make([]byte, file.Size)
	if _, err = io.ReadFull(srcFile, buf); err != nil {
		global.LOG.Error(title, zap.Errors("读取文件", []error{err}))
		httpModule.NewForbidden().SetErrorf("读取文件失败：%w", err).WithAccept(c)
		return
	}

	// 创建 zip reader
	zipReader, err = zip.NewReader(bytes.NewReader(buf), file.Size)
	if err != nil {
		global.LOG.Error(title, zap.Errors("创建zip reader", []error{err}))
		httpModule.NewForbidden().SetErrorf("创建zip容器失败：%w", err).WithAccept(c)
		return
	}

	// 准备写入新 zip 到内存
	var outBuf bytes.Buffer
	zipWriter := zip.NewWriter(&outBuf)

	for _, f := range zipReader.File {
		if f.FileInfo().IsDir() {
			continue
		}

		// 打开原文件
		rc, err = f.Open()
		if err != nil {
			global.LOG.Error(title, zap.Errors("读取压缩包文件", []error{err}))
			httpModule.NewForbidden().SetErrorf("读取压缩包文件失败：%w", err).WithAccept(c)
			return
		}

		// 创建新 zip 内的文件条目
		header := &zip.FileHeader{Name: f.Name, Method: zip.Deflate}
		header.SetModTime(f.ModTime())
		header.SetMode(f.Mode())
		w, err := zipWriter.CreateHeader(header)
		if err != nil {
			_ = rc.Close()
			global.LOG.Error(title, zap.Errors("重建压缩文件", []error{err}))
			httpModule.NewForbidden(httpModule.Errorf("重建压缩文件失败：%w", err)).WithAccept(c)
			return
		}

		// 拷贝内容
		if _, err := io.Copy(w, rc); err != nil {
			_ = rc.Close()
			global.LOG.Error(title, zap.Errors("拷贝压缩文件内容", []error{err}))
			httpModule.NewForbidden().SetErrorf("拷贝压缩文件内容失败：%w", err).WithAccept(c)
			return
		}
		_ = rc.Close()
	}

	// 关闭 zip writer
	if err = zipWriter.Close(); err != nil {
		global.LOG.Error(title, zap.Errors("保存新压缩文件", []error{err}))
		httpModule.NewForbidden().SetErrorf("保存新压缩文件失败：%w", err).WithAccept(c)
		return
	}

	fs := filesystemV2.FileApp.NewByRel(fmt.Sprintf("%s/repacked.zip", global.CONFIG.Rezip.OutDir))
	global.LOG.Info(title, zap.String("保存路径", fs.GetFullPath()))
	if err = os.MkdirAll(fs.GetBasePath(), os.ModePerm); err != nil {
		global.LOG.Error(title, zap.Errors("创建目录", []error{err}))
		httpModule.NewForbidden().SetErrorf("创建目录失败：%w", err).WithAccept(c)
		return
	}

	// 保存到本地文件
	if err = os.WriteFile(fs.GetFullPath(), outBuf.Bytes(), 0644); err != nil {
		global.LOG.Error(title, zap.Errors("保存压缩文件", []error{err}))
		httpModule.NewForbidden().SetErrorf("保存压缩文件失败：%w", err).WithAccept(c)
		return
	}

	httpModule.NewOK().SetMsg("上传成功").SetData(gin.H{"to": c.Request.Host + "/upload/rezip/repacked.zip"}).WithAccept(c)
}
