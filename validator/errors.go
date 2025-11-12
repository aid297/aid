package validator

import (
	"fmt"
	"reflect"

	"github.com/aid297/aid/array"
	"github.com/aid297/aid/myError"
	"github.com/aid297/aid/operation"
)

type (
	ValidateError struct{ myError.MyError }
	RequiredError struct{ myError.MyError }
	EmailError    struct{ myError.MyError }
	TimeError     struct{ myError.MyError }
	LengthError   struct{ myError.MyError }
	RuleError     struct{ myError.MyError }
)

var (
	ValidateErr ValidateError
	RequiredErr RequiredError
	EmailErr    EmailError
	TimeErr     TimeError
	LengthErr   LengthError
	RuleErr     RuleError
)

func (*ValidateError) New(msg string) myError.IMyError {
	return &ValidateError{myError.MyError{Msg: array.NewDestruction("验证错误", msg).JoinWithoutEmpty("：")}}
}
func (my *ValidateError) NewFormat(format string, messages ...any) myError.IMyError {
	return &ValidateError{myError.MyError{Msg: fmt.Sprintf(format, messages...)}}
}

func (*ValidateError) Wrap(err error) myError.IMyError {
	return &ValidateError{myError.MyError{Msg: fmt.Errorf("验证错误"+operation.Ternary(err != nil, "：%w", "%w"), err).Error()}}
}

func (*ValidateError) Panic() myError.IMyError {
	return &ValidateError{MyError: myError.MyError{Msg: "验证错误"}}
}

func (my *ValidateError) Error() string { return my.MyError.Msg }

func (my *ValidateError) Is(target error) bool { return reflect.DeepEqual(target, &ValidateErr) }

func (*RequiredError) New(msg string) myError.IMyError {
	return &RequiredError{myError.MyError{Msg: fmt.Sprintf("[%s]必填", msg)}}
}

func (*RequiredError) Wrap(err error) myError.IMyError {
	return &RequiredError{myError.MyError{Msg: fmt.Errorf("[%w]必填", err).Error()}}
}

func (*RequiredError) Panic() myError.IMyError {
	return &RequiredError{MyError: myError.MyError{Msg: "缺少必填项目"}}
}

func (my *RequiredError) Error() string { return my.MyError.Msg }

func (my *RequiredError) Is(target error) bool { return reflect.DeepEqual(target, &RequiredErr) }

func (*EmailError) New(msg string) myError.IMyError {
	return &EmailError{myError.MyError{Msg: fmt.Sprintf("[%s]不是有效的邮箱格式", msg)}}
}

func (*EmailError) Wrap(err error) myError.IMyError {
	return &EmailError{myError.MyError{Msg: fmt.Errorf("[%w]不是有效的邮箱格式", err).Error()}}
}

func (*EmailError) Panic() myError.IMyError {
	return &EmailError{myError.MyError{Msg: "邮箱格式错误"}}
}

func (my *EmailError) Error() string { return my.MyError.Msg }

func (my *EmailError) Is(target error) bool { return reflect.DeepEqual(target, &EmailErr) }

func (*TimeError) New(msg string) myError.IMyError {
	return &TimeError{myError.MyError{Msg: fmt.Sprintf("[%s]不是有效的邮箱格式", msg)}}
}

func (*TimeError) Wrap(err error) myError.IMyError {
	return &TimeError{myError.MyError{Msg: fmt.Errorf("[%w]不是有效的邮箱格式", err).Error()}}
}

func (*TimeError) Panic() myError.IMyError {
	return &TimeError{myError.MyError{Msg: "时间格式错误"}}
}

func (my *TimeError) Error() string { return my.MyError.Msg }

func (my *TimeError) Is(target error) bool { return reflect.DeepEqual(target, &TimeErr) }

func (my *TimeError) NewFormat(format string, messages ...any) myError.IMyError {
	return &TimeError{myError.MyError{Msg: fmt.Sprintf(format, messages...)}}
}

func (*LengthError) New(msg string) myError.IMyError {
	return &LengthError{myError.MyError{Msg: fmt.Sprintf("[%s]长度错误", msg)}}
}

func (*LengthError) Wrap(err error) myError.IMyError {
	return &LengthError{myError.MyError{Msg: fmt.Errorf("[%w]长度错误", err).Error()}}
}

func (*LengthError) Panic() myError.IMyError {
	return &LengthError{myError.MyError{Msg: "长度错误"}}
}

func (my *LengthError) Error() string { return my.MyError.Msg }

func (my *LengthError) Is(target error) bool { return reflect.DeepEqual(target, &LengthErr) }

func (my *LengthError) NewFormat(format string, messages ...any) myError.IMyError {
	return &LengthError{myError.MyError{Msg: fmt.Sprintf(format, messages...)}}
}

func (*RuleError) New(msg string) myError.IMyError {
	return &RuleError{myError.MyError{Msg: fmt.Sprintf("[%s]规则错误", msg)}}
}

func (*RuleError) Wrap(err error) myError.IMyError {
	return &RuleError{myError.MyError{Msg: fmt.Errorf("[%w]规则错误", err).Error()}}
}

func (*RuleError) Panic() myError.IMyError { return &RuleError{myError.MyError{Msg: "规则错误"}} }

func (my *RuleError) Error() string { return my.MyError.Msg }

func (my *RuleError) Is(target error) bool { return reflect.DeepEqual(target, &LengthErr) }

func (my *RuleError) NewFormat(format string, messages ...any) myError.IMyError {
	return &RuleError{myError.MyError{Msg: fmt.Sprintf(format, messages...)}}
}
