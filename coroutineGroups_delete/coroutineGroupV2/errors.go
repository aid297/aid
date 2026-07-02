package coroutineGroupV2

import "errors"

var (
	ErrBatchInvalid    = errors.New("轮数不能为0")
	ErrCapacityInvalid = errors.New("每轮循环数不能为0")
	ErrEmptyFuncs      = errors.New("没有任何待执行任务")
	ErrTimeout         = errors.New("协程执行超时")
)
