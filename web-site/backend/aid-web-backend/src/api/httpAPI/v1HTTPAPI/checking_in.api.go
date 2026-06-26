package v1HTTPAPI

import (
	"errors"
	"fmt"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid/v5"
	"go.uber.org/zap"

	"github.com/aid297/aid/v2/anySlice"
	"github.com/aid297/aid/v2/consts/volumeInfo"
	"github.com/aid297/aid/v2/excel/excelV3/excelReader"
	"github.com/aid297/aid/v2/filesystem"
	"github.com/aid297/aid/v2/validator"
	"github.com/aid297/aid/v2/web-site/backend/aid-web-backend/src/global"
	"github.com/aid297/aid/v2/web-site/backend/aid-web-backend/src/module/httpModule"
	"github.com/aid297/aid/v2/web-site/backend/aid-web-backend/src/module/httpModule/v1HTTPModule/request"
	"github.com/aid297/aid/v2/web-site/backend/aid-web-backend/src/service/httpService/v1HTTPService"
)

type CheckingInAPI struct{}

// Cal 计算
func (*CheckingInAPI) Cal(c *gin.Context) {
	var (
		title        = "计算考勤"
		err          error
		file         *multipart.FileHeader
		newFilename  string
		saveFile     filesystem.Filesystem
		checker      validator.Checker
		form         request.CheckingInCalRequest
		standardDate *v1HTTPService.StandardDateRet
		everyday     map[string]map[string]v1HTTPService.EverydayRet
		monthly      map[string]map[string]string
	)

	if form, checker = validator.WithGin[request.CheckingInCalRequest](c, func(origin any) (err error) {
		form := origin.(*request.CheckingInCalRequest)
		if repeat := anySlice.NewList(form.ForceWorkdays).Intersection(anySlice.NewList(form.ForceAnnualLeaves)).Intersection(anySlice.NewList(form.Holidays2)).Intersection(anySlice.NewList(form.Holidays3)).RemoveEmpty().ToSlice(); len(repeat) > 0 {
			return errors.New("调休、年薪日、双薪日、三薪日存在重复内容")
		}
		return
	}); checker.Invalid() {
		global.LOG.Error(title, zap.Errors(global.ST_BIND_FORM, checker.Errors()))
		httpModule.NewForbidden(httpModule.Content(checker.Errors()), httpModule.Errorf(global.FE_IVALIDED_FORM, checker.Error())).WithAccept(c)
		return
	}

	if file, err = c.FormFile("file"); err != nil {
		global.LOG.Error(title, zap.Any(global.ST_BIND_FORM, fmt.Sprintf("获取上传文件失败：%s", err)))
		httpModule.NewForbidden(httpModule.Errorf("获取文件失败：%s", err)).WithAccept(c)
		return
	}

	newFilename = fmt.Sprintf("%s.%s", uuid.Must(uuid.NewV6()).String(), file.Filename)
	saveFile = filesystem.NewFile(filesystem.Rel(global.CONFIG.CheckingIn.TmpDir, newFilename))
	defer saveFile.Remove()

	if err = c.SaveUploadedFile(file, saveFile.GetFullPath()); err != nil {
		httpModule.NewForbidden(httpModule.Errorf("保存临时文件失败：%s", err)).WithAccept(c)
		return
	}

	r := excelReader.NewReader(
		excelReader.Filename(saveFile.GetFullPath()),
		excelReader.UnzipXMLSizeLimit(10*volumeInfo.MB),
		excelReader.UnzipSizeLimit(100*volumeInfo.MB),
	)

	if standardDate,
		everyday,
		monthly,
		err =
		v1HTTPService.New.CheckingIn().Cal(r, &form); err != nil {
		httpModule.NewForbidden(httpModule.Error(err)).WithAccept(c)
	}

	httpModule.NewOK(
		httpModule.Msg("上传文件成功"),
		httpModule.Content(gin.H{
			"size":         file.Size,
			"filename":     file.Filename,
			"standardDate": standardDate,
			"everyday":     everyday,
			"monthly":      monthly,
		})).WithAccept(c)
}
