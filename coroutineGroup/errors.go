package coroutineGroup

import "errors"

var (
	ErrBatchInvalid    = errors.New("轮数不能为0")
	ErrCapacityInvalid = errors.New("每轮循环数不能为0")
)
