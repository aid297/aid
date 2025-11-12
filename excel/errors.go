package excel

import (
	"fmt"
	"reflect"

	"github.com/aid297/aid/array"
	"github.com/aid297/aid/myError"
	"github.com/aid297/aid/operation"
)

type (
	SetCellError struct{ myError.MyError }
	ReadError    struct{ myError.MyError }
	WriteError   struct{ myError.MyError }
)

var (
	SetCellErr SetCellError
	ReadErr    ReadError
	WriteErr   WriteError
)

func (*ReadError) New(msg string) myError.IMyError {
	return &ReadError{myError.MyError{Msg: array.NewDestruction("读取数据错误", msg).JoinWithoutEmpty("：")}}
}

func (*ReadError) Wrap(err error) myError.IMyError {
	return &ReadError{myError.MyError{Msg: fmt.Errorf("读取数据错误"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*ReadError) Panic() myError.IMyError {
	return &ReadError{myError.MyError{Msg: "读取数据错误"}}
}

func (my *ReadError) Error() string { return my.Msg }

func (my *ReadError) Is(target error) bool { return reflect.DeepEqual(target, &ReadErr) }

func (*SetCellError) New(msg string) myError.IMyError {
	return &SetCellError{myError.MyError{Msg: array.NewDestruction("设置单元格错误", msg).JoinWithoutEmpty("：")}}
}

func (*SetCellError) Wrap(err error) myError.IMyError {
	return &SetCellError{myError.MyError{Msg: fmt.Errorf("设置单元格错误"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*SetCellError) Panic() myError.IMyError {
	return &SetCellError{myError.MyError{Msg: "设置单元格错误"}}
}

func (my *SetCellError) Error() string { return my.Msg }

func (my *SetCellError) Is(target error) bool { return reflect.DeepEqual(target, &SetCellErr) }

func (*WriteError) New(msg string) myError.IMyError {
	return &WriteError{myError.MyError{Msg: array.NewDestruction("写入数据错误", msg).JoinWithoutEmpty("：")}}
}

func (*WriteError) Wrap(err error) myError.IMyError {
	return &WriteError{myError.MyError{Msg: fmt.Errorf("写入数据错误"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*WriteError) Panic() myError.IMyError {
	return &WriteError{myError.MyError{Msg: "写入数据错误"}}
}

func (my *WriteError) Error() string { return my.Msg }

func (my *WriteError) Is(target error) bool { return reflect.DeepEqual(target, &WriteErr) }
