package validatorV3

import "errors"

var (
	ErrRequired      = errors.New("必填")
	ErrNotEmpty      = errors.New("不能为空")
	ErrInvalidLength = errors.New("长度错误")
	ErrInvalidValue  = errors.New("内容错误")
	ErrInvalidFormat = errors.New("格式错误")
	ErrInvalidType   = errors.New("类型错误")
)
