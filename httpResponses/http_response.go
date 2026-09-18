package httpResponses

import "fmt"

type HTTPResponse[T any] struct {
	HttpCode uint16  `json:"http_code,omitempty"`
	ErrCode  *string `json:"err_code,omitempty"`
	Msg      string  `json:"msg,omitempty"`
	Content  T       `json:"data,omitempty"`
}

func NewOK[T any]() HTTPResponse[T] {
	return HTTPResponse[T]{HttpCode: 200, Msg: "操作成功"}
}

func NewCreated[T any]() HTTPResponse[T] {
	return HTTPResponse[T]{HttpCode: 201, Msg: "创建成功"}
}

func NewUpdated[T any]() HTTPResponse[T] {
	return HTTPResponse[T]{HttpCode: 202, Msg: "编辑成功"}
}

func NewDeleted[T any]() HTTPResponse[T] {
	return HTTPResponse[T]{HttpCode: 204, Msg: "删除成功"}
}

func NewNotFound[T any]() HTTPResponse[T] {
	return HTTPResponse[T]{HttpCode: 404}
}

func NewForbidden[T any]() HTTPResponse[T] {
	return HTTPResponse[T]{HttpCode: 403}
}

func NewNotAuthorization() HTTPResponse[any] {
	return HTTPResponse[any]{HttpCode: 401}
}

func NewBadRequest() HTTPResponse[any] {
	return HTTPResponse[any]{HttpCode: 400}
}

func NewInternalError() HTTPResponse[any] {
	return HTTPResponse[any]{HttpCode: 500}
}

func (my *HTTPResponse[T]) SetMsg(msg string) *HTTPResponse[T] { my.Msg = msg; return my }

func (my *HTTPResponse[T]) SetErrCode(errCode string) *HTTPResponse[T] {
	my.ErrCode = &errCode
	return my
}

func (my *HTTPResponse[T]) SetError(err error) *HTTPResponse[T] { my.Msg = err.Error(); return my }

func (my *HTTPResponse[T]) Format(format string, a ...any) *HTTPResponse[T] {
	my.Msg = fmt.Sprintf(format, a...)
	return my
}

func (my *HTTPResponse[T]) SetContent(content T) *HTTPResponse[T] { my.Content = content; return my }
