package httpModule

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ******************** HTTPResponse 定义 ******************** //
type HTTPResponse struct {
	Code    int    `json:"code" swaggertype:"integer"`
	Msg     string `json:"msg" swaggertype:"string"`
	Content any    `json:"content" swaggertype:"object" x-nullable:"true"`
}

func New(attrs ...HTTPResponseAttributer) HTTPResponse {
	ins := HTTPResponse{}
	for idx := range attrs {
		attrs[idx].Register(&ins)
	}
	return ins
}

func (my HTTPResponse) SetAttrs(attrs ...HTTPResponseAttributer) HTTPResponse {
	for idx := range attrs {
		attrs[idx].Register(&my)
	}
	return my
}

func NewOK(attrs ...HTTPResponseAttributer) HTTPResponse { return New(OK()).SetAttrs(attrs...) }
func NewCreated(attrs ...HTTPResponseAttributer) HTTPResponse {
	return New(Created()).SetAttrs(attrs...)
}
func NewUpdated(attrs ...HTTPResponseAttributer) HTTPResponse {
	return New(Updated()).SetAttrs(attrs...)
}
func NewDeleted(attrs ...HTTPResponseAttributer) HTTPResponse {
	return New(Deleted()).SetAttrs(attrs...)
}
func NewUnauthorized(attrs ...HTTPResponseAttributer) HTTPResponse {
	return New(Unauthorized()).SetAttrs(attrs...)
}
func NewUnPermission(attrs ...HTTPResponseAttributer) HTTPResponse {
	return New(UnPermission()).SetAttrs(attrs...)
}
func NewForbidden(attrs ...HTTPResponseAttributer) HTTPResponse {
	return New(Forbidden()).SetAttrs(attrs...)
}
func NewNotFound(attrs ...HTTPResponseAttributer) HTTPResponse {
	return New(NotFound()).SetAttrs(attrs...)
}
func NewInternalServerError(attrs ...HTTPResponseAttributer) HTTPResponse {
	return New(InternalServerError()).SetAttrs(attrs...)
}

func (my HTTPResponse) SetMsg(msg string) HTTPResponse { return my.SetAttrs(Msg(msg)) }

func (my HTTPResponse) SetError(err error) HTTPResponse { return my.SetAttrs(Error(err)) }

func (my HTTPResponse) SetErrorf(format string, a ...any) HTTPResponse {
	return my.SetAttrs(Errorf(format, a...))
}

func (my HTTPResponse) SetData(data any) HTTPResponse { return my.SetAttrs(Content(data)) }

func (my HTTPResponse) Raw() (int, any) { return my.Code, my }

func (my HTTPResponse) WithAccept(c *gin.Context) {
	switch c.GetHeader("accept") {
	case "application/json":
		my.JSON(c)
	case "application/xml":
		my.XML(c)
	case "application/yaml":
		my.YAML(c)
	case "application/toml":
		my.TOML(c)
	default:
		my.JSON(c)
	}
}

func (my HTTPResponse) JSON(c *gin.Context) { c.JSON(my.Code, my) }
func (my HTTPResponse) YAML(c *gin.Context) { c.YAML(my.Code, my) }
func (my HTTPResponse) TOML(c *gin.Context) { c.TOML(my.Code, my) }
func (my HTTPResponse) XML(c *gin.Context)  { c.XML(my.Code, my) }

func (my HTTPResponse) WithoutCodeJSON(c *gin.Context) { c.JSON(http.StatusOK, my) }
func (my HTTPResponse) WithoutCodeYAML(c *gin.Context) { c.YAML(http.StatusOK, my) }
func (my HTTPResponse) WithoutCodeTOML(c *gin.Context) { c.TOML(http.StatusOK, my) }
func (my HTTPResponse) WithoutCodeXML(c *gin.Context)  { c.XML(http.StatusOK, my) }

// ******************** HTTPResponseAttributer 实现 ******************** //
type (
	HTTPResponseAttributer interface{ Register(res *HTTPResponse) }

	AttrCode struct {
		code int
		msg  string
	}
	AttrMsg  struct{ msg string }
	AttrData struct{ data any }
)

func OK() HTTPResponseAttributer      { return AttrCode{code: http.StatusOK, msg: "操作成功"} }
func Created() HTTPResponseAttributer { return AttrCode{code: http.StatusCreated, msg: "新建成功"} }
func Updated() HTTPResponseAttributer {
	return AttrCode{code: http.StatusAccepted, msg: "更新成功"}
}
func Deleted() HTTPResponseAttributer {
	return AttrCode{code: http.StatusNoContent, msg: "删除成功"}
}
func Unauthorized() HTTPResponseAttributer {
	return AttrCode{code: http.StatusUnauthorized, msg: "未登录"}
}
func UnPermission() HTTPResponseAttributer {
	return AttrCode{code: http.StatusNotAcceptable, msg: "无权操作"}
}
func Forbidden() HTTPResponseAttributer {
	return AttrCode{code: http.StatusForbidden, msg: "操作失败"}
}
func NotFound() HTTPResponseAttributer {
	return AttrCode{code: http.StatusNotFound, msg: "资源不存在"}
}
func InternalServerError() HTTPResponseAttributer {
	return AttrCode{code: http.StatusInternalServerError, msg: "服务器内部错误"}
}
func (my AttrCode) Register(res *HTTPResponse) { res.Code = my.code; res.Msg = my.msg }

func Msg(msg string) HTTPResponseAttributer                 { return AttrMsg{msg: msg} }
func Error(err error) HTTPResponseAttributer                { return AttrMsg{msg: err.Error()} }
func Errorf(format string, a ...any) HTTPResponseAttributer { return Error(fmt.Errorf(format, a...)) }
func (my AttrMsg) Register(res *HTTPResponse)               { res.Msg = my.msg }

func Content(data any) HTTPResponseAttributer  { return AttrData{data: data} }
func (my AttrData) Register(res *HTTPResponse) { res.Content = my.data }
