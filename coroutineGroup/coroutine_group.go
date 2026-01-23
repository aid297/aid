package coroutineGroup

import (
	"sync"

	"github.com/aid297/aid/operation/operationV2"
)

type (
	CoroutineGroup[T any] struct {
		Error      error
		OK         bool
		Results    []*Result[T]
		batches    uint
		capacities uint
		sw         sync.WaitGroup
	}

	Result[T any] struct {
		Data   T
		Error  error
		IsSkip bool
	}
)

// NewCoroutineGroup 创建协程组实例
func NewCoroutineGroup[T any]() *CoroutineGroup[T] {
	return &CoroutineGroup[T]{sw: sync.WaitGroup{}, OK: true}
}

// GetBatches 计算批次数
func GetBatches(total, capacities int) uint { return uint((total + capacities - 1) / capacities) }

// SetBatches 设置批次数
func (my *CoroutineGroup[T]) SetBatches(batches uint) *CoroutineGroup[T] {
	my.batches = batches
	return my
}

// SetCapacity 设置每批次容量
func (my *CoroutineGroup[T]) SetCapacity(capacities uint) *CoroutineGroup[T] {
	my.capacities = capacities
	return my
}

// SetBatchesByCapacities 根据总数和每批次容量计算批次数并设置
func (my *CoroutineGroup[T]) SetBatchesByCapacities(total, capacities int) *CoroutineGroup[T] {
	my.batches = operationV2.NewTernary(operationV2.TrueFn(func() uint { return GetBatches(total, capacities) }), operationV2.FalseValue[uint](1)).GetByValue(total > capacities)
	my.capacities = uint(capacities)
	return my
}

// check 检查参数
func (my *CoroutineGroup[T]) check() error {
	if my.batches == 0 {
		return ErrBatchInvalid
	}
	if my.capacities == 0 {
		return ErrCapacityInvalid
	}

	return nil
}

// GO 执行协程组
func (my *CoroutineGroup[T]) GO(fn func(batch, capacity uint) (result *Result[T])) *CoroutineGroup[T] {
	if err := my.check(); err != nil {
		my.Error = err
		my.OK = false
		return my
	}

	my.Results = make([]*Result[T], 0, my.batches+my.capacities)
	for batch := range my.batches {
		for capacity := range my.capacities {
			my.sw.Add(1)

			go func(b, c uint) {
				defer my.sw.Done()
				var r *Result[T] = fn(b, c)
				my.Results = append(my.Results, r)
				if r.Error != nil {
					my.OK = false
				}
			}(batch, capacity)
		}
		my.sw.Wait()
	}

	return my
}
