package coroutineGroup

import (
	"sync"
)

type (
	CoroutineGroup[T any] struct {
		Error      error
		batches    uint
		capacities uint
		Results    []*Result[T]
		sw         sync.WaitGroup
		OK         bool
	}
	Result[T any] struct {
		Data  T
		Error error
	}
)

func New[T any]() *CoroutineGroup[T] { return &CoroutineGroup[T]{sw: sync.WaitGroup{}, OK: true} }

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

func (my *CoroutineGroup[T]) Run(fn func(batch, capacity uint) *Result[T]) *CoroutineGroup[T] {
	if err := my.check(); err != nil {
		my.Error = err
		my.OK = false
		return my
	}

	my.Results = make([]*Result[T], 0, my.batches+my.capacities)
	for batch := range my.batches {
		for capacity := range my.capacities {
			my.sw.Add(1)

			var (
				r  *Result[T]
				ch = make(chan struct{})
			)

			go func(b, c uint) {
				defer my.sw.Done()
				r = fn(b, c)
				ch <- struct{}{}
			}(batch, capacity)

			<-ch
			my.Results = append(my.Results, r)
			if r.Error != nil {
				my.OK = false
			}
		}
		my.sw.Wait()
	}

	return my
}

func (my *CoroutineGroup[T]) RunUntilError(fn func(batch, capacity uint) *Result[T]) *CoroutineGroup[T] {
	if err := my.check(); err != nil {
		my.Error = err
		my.OK = false
		return my
	}

	my.Results = make([]*Result[T], 0, my.batches+my.capacities)
	for batch := range my.batches {
		for capacity := range my.capacities {
			my.sw.Add(1)

			var (
				r  *Result[T]
				ch = make(chan struct{})
			)

			go func(b, c uint) {
				defer my.sw.Done()
				r = fn(b, c)
				ch <- struct{}{}
			}(batch, capacity)

			<-ch
			my.Results = append(my.Results, r)
			if r.Error != nil {
				my.OK = false
				return my
			}
		}
		my.sw.Wait()
	}

	return my
}
