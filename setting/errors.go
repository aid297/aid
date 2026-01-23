package setting

import "errors"

var (
	ErrConfigNotSet       = errors.New("配置文件路径未设置")
	ErrConfigFileNotFound = errors.New("配置文件未找到")
	ErrConfigReadFailed   = errors.New("读取配置文件失败")
)
