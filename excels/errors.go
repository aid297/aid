package excels

import (
	"reflect"

	"github.com/aid297/aid/v3/anySlices"
	"github.com/aid297/aid/v3/myError"
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
	return &ReadError{myError.MyError{Msg: anySlices.New(anySlices.Items("读取数据错误", msg)).JoinNotEmpty("：")}}
}

func (*ReadError) Wrap(err error) myError.IMyError {
	return &ReadError{myError.MyError{Msg: anySlices.New(anySlices.Items("读取数据错误", err.Error())).JoinNotEmpty("：")}}
}

func (*ReadError) Panic() myError.IMyError {
	return &ReadError{myError.MyError{Msg: "读取数据错误"}}
}

func (my *ReadError) Error() string { return my.Msg }

func (my *ReadError) Is(target error) bool { return reflect.DeepEqual(target, &ReadErr) }

func (*SetCellError) New(msg string) myError.IMyError {
	return &SetCellError{myError.MyError{Msg: anySlices.New(anySlices.Items("设置单元格错误", msg)).JoinNotEmpty("：")}}
}

func (*SetCellError) Wrap(err error) myError.IMyError {
	return &SetCellError{myError.MyError{Msg: anySlices.NewItems("设置单元格错误", err.Error()).JoinNotEmpty("：")}}
}

func (*SetCellError) Panic() myError.IMyError {
	return &SetCellError{myError.MyError{Msg: "设置单元格错误"}}
}

func (my *SetCellError) Error() string { return my.Msg }

func (my *SetCellError) Is(target error) bool { return reflect.DeepEqual(target, &SetCellErr) }

func (*WriteError) New(msg string) myError.IMyError {
	return &WriteError{myError.MyError{Msg: anySlices.New(anySlices.Items("写入数据错误", msg)).JoinNotEmpty("：")}}
}

func (*WriteError) Wrap(err error) myError.IMyError {
	return &WriteError{myError.MyError{Msg: anySlices.New(anySlices.Items("写入数据错误", err.Error())).JoinNotEmpty("：")}}
}

func (*WriteError) Panic() myError.IMyError {
	return &WriteError{myError.MyError{Msg: "写入数据错误"}}
}

func (my *WriteError) Error() string { return my.Msg }

func (my *WriteError) Is(target error) bool { return reflect.DeepEqual(target, &WriteErr) }
