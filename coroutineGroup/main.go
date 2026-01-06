package coroutineGroup

import (
	"sync"
)

type (
	CoroutineGroup[T any] struct {
		Error    error
		batch    uint
		capacity uint
		Results  []*Result[T]
		sw       sync.WaitGroup
		OK       bool
	}
	Result[T any] struct {
		Data  T
		Error error
	}
)

func New[T any]() *CoroutineGroup[T] { return &CoroutineGroup[T]{sw: sync.WaitGroup{}} }

func (my *CoroutineGroup[T]) SetBatches(batches uint) *CoroutineGroup[T] {
	my.batch = batches
	return my
}

func (my *CoroutineGroup[T]) SetCapacity(capacity uint) *CoroutineGroup[T] {
	my.capacity = capacity
	return my
}

func (my *CoroutineGroup[T]) check() error {
	if my.batch == 0 {
		return ErrBatchInvalid
	}
	if my.capacity == 0 {
		return ErrCapacityInvalid
	}

	return nil
}

func (my *CoroutineGroup[T]) Run(fn func() *Result[T]) *CoroutineGroup[T] {
	my.Results = make([]*Result[T], 0, my.batch+my.capacity)
	for range my.batch {
		for range my.capacity {
			my.sw.Add(1)
			defer my.sw.Done()
			r := fn()
			if r.Error != nil {
				my.OK = false
			}
			my.Results = append(my.Results, fn())
		}
		my.sw.Wait()
	}

	return my
}

func (my *CoroutineGroup[T]) RunUntilError(fn func() *Result[T]) *CoroutineGroup[T] {
	my.Results = make([]*Result[T], 0, my.batch+my.capacity)
	for range my.batch {
		for range my.capacity {
			my.sw.Add(1)
			defer my.sw.Done()
			r := fn()
			my.Results = append(my.Results, fn())
			if r.Error != nil {
				my.OK = false
				return my
			}
		}
		my.sw.Wait()
	}

	return my
}
