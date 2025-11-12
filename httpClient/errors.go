package httpClient

import (
	"fmt"
	"reflect"

	"github.com/aid297/aid/array"
	"github.com/aid297/aid/myError"
	"github.com/aid297/aid/operation"
)

type (
	ReadResponseError    struct{ myError.MyError }
	UrlEmptyError        struct{ myError.MyError }
	GenerateCertError    struct{ myError.MyError }
	GenerateRequestError struct{ myError.MyError }
	UnmarshalXmlError    struct{ myError.MyError }
	UnmarshalJsonError   struct{ myError.MyError }
	SetSteamBodyError    struct{ myError.MyError }
	SetFormBodyError     struct{ myError.MyError }
	SetXmlBodyError      struct{ myError.MyError }
	SetJsonBodyError     struct{ myError.MyError }
	WriteResponseError   struct{ myError.MyError }
)

var (
	ReadResponseErr    ReadResponseError
	UrlEmptyErr        UrlEmptyError
	GenerateCertErr    GenerateCertError
	GenerateRequestErr GenerateRequestError
	UnmarshalXmlErr    UnmarshalXmlError
	UnmarshalJsonErr   UnmarshalJsonError
	SetSteamBodyErr    SetSteamBodyError
	SetFormBodyErr     SetFormBodyError
	SetXmlBodyErr      SetXmlBodyError
	SetJsonBodyErr     SetJsonBodyError
	WriteResponseErr   WriteResponseError
)

func (*ReadResponseError) New(msg string) myError.IMyError {
	return &ReadResponseError{MyError: myError.MyError{Msg: array.New([]string{"读取响应体失败", msg}).JoinWithoutEmpty("：")}}
}

func (*ReadResponseError) Wrap(err error) myError.IMyError {
	return &ReadResponseError{MyError: myError.MyError{Msg: fmt.Errorf("读取响应体失败"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*ReadResponseError) Panic() myError.IMyError {
	return &ReadResponseError{MyError: myError.MyError{Msg: "读取响应体失败"}}
}

func (my *ReadResponseError) Error() string { return my.MyError.Msg }

func (my *ReadResponseError) Is(target error) bool {
	return reflect.DeepEqual(target, &ReadResponseErr)
}

func (*UrlEmptyError) New(msg string) myError.IMyError {
	return &UrlEmptyError{MyError: myError.MyError{Msg: array.New([]string{"url不能为空", msg}).JoinWithoutEmpty("：")}}
}

func (*UrlEmptyError) Wrap(err error) myError.IMyError {
	return &UrlEmptyError{MyError: myError.MyError{Msg: fmt.Errorf("url不能为空"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*UrlEmptyError) Panic() myError.IMyError {
	return &UrlEmptyError{MyError: myError.MyError{Msg: "url不能为空"}}
}

func (my *UrlEmptyError) Error() string { return my.MyError.Msg }

func (my *UrlEmptyError) Is(target error) bool {
	return reflect.DeepEqual(target, &UrlEmptyErr)
}

func (*GenerateCertError) New(msg string) myError.IMyError {
	return &GenerateCertError{MyError: myError.MyError{Msg: array.New([]string{"生成证书失败", msg}).JoinWithoutEmpty("：")}}
}

func (*GenerateCertError) Wrap(err error) myError.IMyError {
	return &GenerateCertError{MyError: myError.MyError{Msg: fmt.Errorf("生成证书失败"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*GenerateCertError) Panic() myError.IMyError {
	return &GenerateCertError{MyError: myError.MyError{Msg: "生成证书失败"}}
}

func (my *GenerateCertError) Error() string { return my.MyError.Msg }

func (my *GenerateCertError) Is(target error) bool {
	return reflect.DeepEqual(target, &GenerateCertErr)
}

func (*GenerateRequestError) New(msg string) myError.IMyError {
	return &GenerateRequestError{MyError: myError.MyError{Msg: array.New([]string{"生成请求对象失败", msg}).JoinWithoutEmpty("：")}}
}

func (*GenerateRequestError) Wrap(err error) myError.IMyError {
	return &GenerateRequestError{MyError: myError.MyError{Msg: fmt.Errorf("生成请求对象失败"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*GenerateRequestError) Panic() myError.IMyError {
	return &GenerateRequestError{MyError: myError.MyError{Msg: "生成请求对象失败"}}
}

func (my *GenerateRequestError) Error() string { return my.MyError.Msg }

func (my *GenerateRequestError) Is(target error) bool {
	return reflect.DeepEqual(target, &GenerateRequestErr)
}

func (*UnmarshalXmlError) New(msg string) myError.IMyError {
	return &UnmarshalXmlError{MyError: myError.MyError{Msg: array.New([]string{"获取xml格式响应体错误", msg}).JoinWithoutEmpty("：")}}
}

func (*UnmarshalXmlError) Wrap(err error) myError.IMyError {
	return &UnmarshalXmlError{MyError: myError.MyError{Msg: fmt.Errorf("获取xml格式响应体错误"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*UnmarshalXmlError) Panic() myError.IMyError {
	return &UnmarshalXmlError{MyError: myError.MyError{Msg: "获取xml格式响应体错误"}}
}

func (my *UnmarshalXmlError) Error() string { return my.MyError.Msg }

func (my *UnmarshalXmlError) Is(target error) bool {
	return reflect.DeepEqual(target, &UnmarshalXmlErr)
}

func (*UnmarshalJsonError) New(msg string) myError.IMyError {
	return &UnmarshalJsonError{MyError: myError.MyError{Msg: array.New([]string{"获取json格式响应体错误", msg}).JoinWithoutEmpty("：")}}
}

func (*UnmarshalJsonError) Wrap(err error) myError.IMyError {
	return &UnmarshalJsonError{MyError: myError.MyError{Msg: fmt.Errorf("获取json格式响应体错误"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*UnmarshalJsonError) Panic() myError.IMyError {
	return &UnmarshalJsonError{MyError: myError.MyError{Msg: "获取json格式响应体错误"}}
}

func (my *UnmarshalJsonError) Error() string { return my.MyError.Msg }

func (my *UnmarshalJsonError) Is(target error) bool {
	return reflect.DeepEqual(target, &UnmarshalJsonErr)
}

func (*SetSteamBodyError) New(msg string) myError.IMyError {
	return &SetSteamBodyError{MyError: myError.MyError{Msg: array.New([]string{"设置二进制请求体失败", msg}).JoinWithoutEmpty("：")}}
}

func (*SetSteamBodyError) Wrap(err error) myError.IMyError {
	return &SetSteamBodyError{MyError: myError.MyError{Msg: fmt.Errorf("设置二进制请求体失败"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*SetSteamBodyError) Panic() myError.IMyError {
	return &SetSteamBodyError{MyError: myError.MyError{Msg: "设置二进制请求体失败"}}
}

func (my *SetSteamBodyError) Error() string { return my.MyError.Msg }

func (my *SetSteamBodyError) Is(target error) bool {
	return reflect.DeepEqual(target, &SetSteamBodyErr)
}

func (*SetFormBodyError) New(msg string) myError.IMyError {
	return &SetFormBodyError{MyError: myError.MyError{Msg: array.New([]string{"设置表单数据请求体失败", msg}).JoinWithoutEmpty("：")}}
}

func (*SetFormBodyError) Wrap(err error) myError.IMyError {
	return &SetFormBodyError{MyError: myError.MyError{Msg: fmt.Errorf("设置表单数据请求体失败"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*SetFormBodyError) Panic() myError.IMyError {
	return &SetFormBodyError{MyError: myError.MyError{Msg: "设置表单数据请求体失败"}}
}

func (my *SetFormBodyError) Error() string { return my.MyError.Msg }

func (my *SetFormBodyError) Is(target error) bool { return reflect.DeepEqual(target, &SetFormBodyErr) }

func (*SetXmlBodyError) New(msg string) myError.IMyError {
	return &SetXmlBodyError{MyError: myError.MyError{Msg: array.New([]string{"设置xml请求体失败", msg}).JoinWithoutEmpty("：")}}
}

func (*SetXmlBodyError) Wrap(err error) myError.IMyError {
	return &SetXmlBodyError{MyError: myError.MyError{Msg: fmt.Errorf("设置xml请求体失败"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*SetXmlBodyError) Panic() myError.IMyError {
	return &SetXmlBodyError{MyError: myError.MyError{Msg: "设置xml请求体失败"}}
}

func (my *SetXmlBodyError) Error() string { return my.MyError.Msg }

func (my *SetXmlBodyError) Is(target error) bool { return reflect.DeepEqual(target, &SetXmlBodyErr) }

func (*SetJsonBodyError) New(msg string) myError.IMyError {
	return &SetJsonBodyError{MyError: myError.MyError{Msg: array.New([]string{"设置json请求体失败", msg}).JoinWithoutEmpty("：")}}
}

func (*SetJsonBodyError) Wrap(err error) myError.IMyError {
	return &SetJsonBodyError{MyError: myError.MyError{Msg: fmt.Errorf("设置json请求体失败"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*SetJsonBodyError) Panic() myError.IMyError {
	return &SetJsonBodyError{MyError: myError.MyError{Msg: "设置json请求体失败"}}
}

func (my *SetJsonBodyError) Error() string { return my.MyError.Msg }

func (my *SetJsonBodyError) Is(target error) bool { return reflect.DeepEqual(target, &SetJsonBodyErr) }

func (*WriteResponseError) New(msg string) myError.IMyError {
	return &WriteResponseError{MyError: myError.MyError{Msg: array.New([]string{"写入响应失败", msg}).JoinWithoutEmpty("：")}}
}

func (*WriteResponseError) Wrap(err error) myError.IMyError {
	return &WriteResponseError{MyError: myError.MyError{Msg: fmt.Errorf("写入响应失败"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*WriteResponseError) Panic() myError.IMyError {
	return &WriteResponseError{MyError: myError.MyError{Msg: "写入响应失败"}}
}

func (my *WriteResponseError) Error() string { return my.MyError.Msg }

func (my *WriteResponseError) Is(target error) bool {
	return reflect.DeepEqual(target, &WriteResponseErr)
}
