package httpModule

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
func NewUnprocessableEntity(attrs ...HTTPResponseAttributer) HTTPResponse {
	return New(UnprocessableEntity()).SetAttrs(attrs...)
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

type httpBodyFormat string

const (
	httpBodyFormatJSON httpBodyFormat = "json"
	httpBodyFormatXML  httpBodyFormat = "xml"
	httpBodyFormatYAML httpBodyFormat = "yaml"
	httpBodyFormatTOML httpBodyFormat = "toml"
)

type acceptCandidate struct {
	value string
	q     float64
	pos   int
}

func parseAcceptHeader(header string) []acceptCandidate {
	header = strings.TrimSpace(header)
	if header == "" {
		return nil
	}
	parts := strings.Split(header, ",")
	out := make([]acceptCandidate, 0, len(parts))
	for i, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		mediaRange := part
		q := 1.0
		if semi := strings.Index(part, ";"); semi >= 0 {
			mediaRange = strings.TrimSpace(part[:semi])
			params := strings.Split(part[semi+1:], ";")
			for _, param := range params {
				param = strings.TrimSpace(param)
				if param == "" {
					continue
				}
				kv := strings.SplitN(param, "=", 2)
				if len(kv) != 2 {
					continue
				}
				if strings.EqualFold(strings.TrimSpace(kv[0]), "q") {
					if parsed, err := strconv.ParseFloat(strings.TrimSpace(kv[1]), 64); err == nil {
						q = parsed
					}
				}
			}
		}
		if mediaRange == "" {
			continue
		}
		out = append(out, acceptCandidate{value: strings.ToLower(mediaRange), q: q, pos: i})
	}
	return out
}

func matchAcceptFormat(candidate string) (httpBodyFormat, bool) {
	candidate = strings.TrimSpace(strings.ToLower(candidate))
	if candidate == "" {
		return "", false
	}
	if candidate == "*/*" || candidate == "application/*" {
		return httpBodyFormatJSON, true
	}
	if strings.HasSuffix(candidate, "+json") || candidate == "application/json" || candidate == "text/json" {
		return httpBodyFormatJSON, true
	}
	if strings.HasSuffix(candidate, "+xml") || candidate == "application/xml" || candidate == "text/xml" {
		return httpBodyFormatXML, true
	}
	if candidate == "application/yaml" || candidate == "text/yaml" || candidate == "application/x-yaml" || candidate == "text/x-yaml" {
		return httpBodyFormatYAML, true
	}
	if candidate == "application/toml" || candidate == "text/toml" {
		return httpBodyFormatTOML, true
	}
	return "", false
}

func chooseResponseFormat(acceptHeader string) httpBodyFormat {
	best := acceptCandidate{q: -1}
	bestFormat := httpBodyFormatJSON
	for _, candidate := range parseAcceptHeader(acceptHeader) {
		format, ok := matchAcceptFormat(candidate.value)
		if !ok {
			continue
		}
		if candidate.q > best.q || (candidate.q == best.q && candidate.pos < best.pos) {
			best = candidate
			bestFormat = format
		}
	}
	return bestFormat
}

func (my HTTPResponse) WithAccept(c *gin.Context) {
	c.Header("Vary", "Accept")
	switch chooseResponseFormat(c.GetHeader("Accept")) {
	case httpBodyFormatXML:
		my.XML(c)
	case httpBodyFormatYAML:
		my.YAML(c)
	case httpBodyFormatTOML:
		my.TOML(c)
	default:
		my.JSON(c)
	}
}

func (my HTTPResponse) WithAcceptWithoutCode(c *gin.Context) {
	c.Header("Vary", "Accept")
	switch chooseResponseFormat(c.GetHeader("Accept")) {
	case httpBodyFormatXML:
		my.WithoutCodeXML(c)
	case httpBodyFormatYAML:
		my.WithoutCodeYAML(c)
	case httpBodyFormatTOML:
		my.WithoutCodeTOML(c)
	default:
		my.WithoutCodeJSON(c)
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
func Created() HTTPResponseAttributer { return AttrCode{code: http.StatusCreated, msg: "创建成功"} }
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
func UnprocessableEntity() HTTPResponseAttributer {
	return AttrCode{code: http.StatusUnprocessableEntity, msg: "表单验证失败"}
}
func NotFound() HTTPResponseAttributer {
	return AttrCode{code: http.StatusNotFound, msg: "资源不存在"}
}
func InternalServerError() HTTPResponseAttributer {
	return AttrCode{code: http.StatusInternalServerError, msg: "其他错误"}
}
func (my AttrCode) Register(res *HTTPResponse) { res.Code = my.code; res.Msg = my.msg }

func Msg(msg string) HTTPResponseAttributer                 { return AttrMsg{msg: msg} }
func Error(err error) HTTPResponseAttributer                { return AttrMsg{msg: err.Error()} }
func Errorf(format string, a ...any) HTTPResponseAttributer { return Error(fmt.Errorf(format, a...)) }
func (my AttrMsg) Register(res *HTTPResponse)               { res.Msg = my.msg }

func Content(data any) HTTPResponseAttributer  { return AttrData{data: data} }
func (my AttrData) Register(res *HTTPResponse) { res.Content = my.data }
