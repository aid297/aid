package coroutineGroup

import (
	"sync"
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

func New[T any]() *CoroutineGroup[T]        { return &CoroutineGroup[T]{sw: sync.WaitGroup{}, OK: true} }
func GetBatches(total, capacities int) uint { return uint((total + capacities - 1) / capacities) }

func (my *CoroutineGroup[T]) SetBatches(batches uint) *CoroutineGroup[T] {
	my.batches = batches
	return my
}

func (my *CoroutineGroup[T]) SetCapacity(capacity uint) *CoroutineGroup[T] {
	my.capacities = capacity
	return my
}

func (my *CoroutineGroup[T]) check() error {
	if my.batches == 0 {
		return ErrBatchInvalid
	}
	if my.capacities == 0 {
		return ErrCapacityInvalid
	}

	return nil
}

func (my *CoroutineGroup[T]) Run(fn func(batch, capacity uint) (result *Result[T])) *CoroutineGroup[T] {
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
