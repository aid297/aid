package validatorV2

import "errors"

var (
	ErrRequired    = errors.New("必填")
	ErrNoZero      = errors.New("不能为空")
	ErrEq          = errors.New("必须等于")
	ErrNotEq       = errors.New("不能等于")
	ErrIn          = errors.New("仅允许")
	ErrNotIn       = errors.New("不允许")
	ErrRegex       = errors.New("内容不匹配")
	ErrParseData   = errors.New("解析数据错误")
	ErrInvalidType = errors.New("无效的类型")

	TooShort = "长度不能小于"
	TooLong  = "长度不能大于"
	TooSmall = "不能小于"
	TooBig   = "不能大于"
)
