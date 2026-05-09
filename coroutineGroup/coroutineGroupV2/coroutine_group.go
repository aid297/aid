package coroutineGroupV2

import (
	"sync"
)

type (
	CoroutineGroup[T any] struct {
		sem   chan struct{}
		funcs []Func[T]
	}

	Func[T any] func() Result[T]

	Result[T any] struct {
		Value T
		Error error
	}
)

func New[T any](limit uint16) *CoroutineGroup[T] {
	if limit == 0 {
		limit = 4
	}

	return &CoroutineGroup[T]{sem: make(chan struct{}, limit)}
}

func (my *CoroutineGroup[T]) SetFunc(funcs ...Func[T]) *CoroutineGroup[T] {
	my.funcs = append(my.funcs, funcs...)
	return my
}

// GO 批量执行
func (my *CoroutineGroup[T]) GO(funcs ...Func[T]) []Result[T] {
	var (
		wg      = sync.WaitGroup{}
		results = make([]Result[T], len(funcs))
	)

	if len(funcs) > 0 {
		my.funcs = funcs
	}

	if len(my.funcs) == 0 {
		return results
	}

	for idx := range my.funcs {
		wg.Add(1)
		my.sem <- struct{}{}
		go func(idx int) { defer wg.Done(); defer func() { <-my.sem }(); results[idx] = funcs[idx]() }(idx)
	}
	wg.Wait()

	return results
}
