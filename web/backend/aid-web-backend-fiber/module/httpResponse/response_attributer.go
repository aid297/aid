package httpResponse

import (
	"fmt"
	"net/http"
)

type (
	ResponseAttributer interface {
		Register(res *response)
	}

	AttrCode struct{ code int }
	AttrMsg  struct{ msg string }
	AttrData struct{ data any }
)

func OK() ResponseAttributer           { return AttrCode{code: http.StatusOK} }
func Created() ResponseAttributer      { return AttrCode{code: http.StatusCreated} }
func Updated() ResponseAttributer      { return AttrCode{code: http.StatusAccepted} }
func Deleted() ResponseAttributer      { return AttrCode{code: http.StatusNoContent} }
func NoLogin() ResponseAttributer      { return AttrCode{code: http.StatusUnauthorized} }
func NoPermission() ResponseAttributer { return AttrCode{code: http.StatusNotAcceptable} }
func Forbidden() ResponseAttributer    { return AttrCode{code: http.StatusForbidden} }
func NotFound() ResponseAttributer     { return AttrCode{code: http.StatusNotFound} }
func (my AttrCode) Register(res *response) {
	res.code = my.code

	switch my.code {
	case http.StatusOK:
		res.msg = "OK"
	case http.StatusCreated:
		res.msg = "新建成功"
	case http.StatusAccepted:
		res.msg = "操作成功"
	case http.StatusNoContent:
		res.msg = "删除成功"
	case http.StatusUnauthorized:
		res.msg = "未登录"
	case http.StatusNotAcceptable:
		res.msg = "无权限"
	case http.StatusForbidden:
		res.msg = "操作失败"
	case http.StatusNotFound:
		res.msg = "资源不存在"
	}
}

func Msg(msg string) ResponseAttributer  { return AttrMsg{msg: msg} }
func Error(err error) ResponseAttributer { return AttrMsg{msg: err.Error()} }
func Errorf(format string, a ...any) ResponseAttributer {
	return AttrMsg{msg: fmt.Errorf(format, a...).Error()}
}
func (my AttrMsg) Register(res *response) { res.msg = my.msg }

func Data(data any) ResponseAttributer     { return AttrData{data: data} }
func (my AttrData) Register(res *response) { res.data = my.data }
