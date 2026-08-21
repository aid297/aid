package coroutineGroups

import (
	"sync"
	"time"

	"github.com/aid297/aid/v3/anySlices"
	"github.com/aid297/aid/v3/operations"
)

type (
	CoroutineGrouper[T any] interface {
		New(limit uint16) CoroutineGrouper[T]
		SetFunc(funcs ...Func[T]) CoroutineGrouper[T]
		SetRetry(opts ...CoroutineGroupRetryAttr) CoroutineGrouper[T]
		GO(funcs ...Func[T]) anySlices.AnySlicer[Result[T]]
		GOBatch(total, capacities int, fn func(batch, capacity uint) (result Result[T])) (anySlices.AnySlicer[Result[T]], error)
	}

	CoroutineGroupImpl[T any] struct {
		sem        chan struct{}
		funcs      []Func[T]
		batches    uint
		capacities uint
		retry      *RetryConfig
	}

	Func[T any] func() Result[T]

	Result[T any] interface {
		IsOK() bool
		IsSkip() bool
		GetError() error
		GetValue() T
		SetOK(ok bool)
		SetSkip(skip bool)
	}

	ResultImpl[T any] struct {
		Error error
		Skip  bool
		Value T
	}
)

func New[T any](limit uint16, attrs ...CoroutineGroupRetryAttr) CoroutineGrouper[T] {
	if limit == 0 {
		limit = 4
	}

	ins := &CoroutineGroupImpl[T]{sem: make(chan struct{}, limit), retry: nil}

	if len(attrs) > 0 {
		ins.retry = &RetryConfig{maxAttempts: 1}
		for idx := range attrs {
			attrs[idx](ins.retry)
		}
	}

	return ins
}

func (my *CoroutineGroupImpl[T]) New(limit uint16) CoroutineGrouper[T] { return New[T](limit) }

func (my *CoroutineGroupImpl[T]) SetFunc(funcs ...Func[T]) CoroutineGrouper[T] {
	my.funcs = append(my.funcs, funcs...)
	return my
}

// SetRetry 设置超时重试配置，每个协程独立享有超时与重试
func (my *CoroutineGroupImpl[T]) SetRetry(attrs ...CoroutineGroupRetryAttr) CoroutineGrouper[T] {
	my.retry = &RetryConfig{maxAttempts: 1}
	for idx := range attrs {
		attrs[idx](my.retry)
	}

	return my
}

// executeWithRetry 执行函数，支持超时与重试，每个协程独立
func (my *CoroutineGroupImpl[T]) executeWithRetry(fn Func[T]) Result[T] {
	// 未配置重试，直接执行
	if my.retry == nil || my.retry.maxAttempts <= 1 {
		return fn()
	}

	cfg := my.retry
	var lastResult Result[T] = &ResultImpl[T]{Error: ErrTimeout}

	for attempt := 0; attempt < cfg.maxAttempts; attempt++ {
		// 非首次执行前等待退避时间
		if attempt > 0 {
			time.Sleep(cfg.calculateWait(attempt))
		}

		// 带超时执行
		result, timedOut := my.executeWithTimeout(fn, cfg.timeout)
		if timedOut {
			lastResult = &ResultImpl[T]{Error: ErrTimeout}
			continue
		}

		// 成功或跳过，无需重试
		if result.IsOK() || result.IsSkip() {
			return result
		}

		lastResult = result
	}

	return lastResult
}

// executeWithTimeout 带超时执行函数，超时返回 (nil, true)
func (my *CoroutineGroupImpl[T]) executeWithTimeout(fn Func[T], timeout time.Duration) (Result[T], bool) {
	if timeout <= 0 {
		return fn(), false
	}

	done := make(chan Result[T], 1)
	go func() { done <- fn() }()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case r := <-done:
		return r, false
	case <-timer.C:
		return nil, true
	}
}

// GO 批量执行
func (my *CoroutineGroupImpl[T]) GO(funcs ...Func[T]) anySlices.AnySlicer[Result[T]] {
	var (
		wg      = sync.WaitGroup{}
		results = anySlices.New(anySlices.Cap[Result[T]](len(funcs)))
	)

	if len(my.funcs) == 0 {
		return results
	}
	my.funcs = funcs

	for idx := range my.funcs {
		wg.Add(1)
		my.sem <- struct{}{}
		go func(idx int) {
			defer wg.Done()
			defer func() { <-my.sem }()
			results.Lock()
			results.Append(my.executeWithRetry(funcs[idx]))
			results.Unlock()
		}(idx)
	}
	wg.Wait()

	return results
}

func (my *CoroutineGroupImpl[T]) GOBatch(total, capacities int, fn func(batch, capacity uint) (result Result[T])) (anySlices.AnySlicer[Result[T]], error) {
	if total == 0 {
		return nil, ErrEmptyFuncs
	}

	my.batches = operations.NewTernary(operations.TrueFn(func() uint { return uint((total + capacities - 1) / capacities) }), operations.FalseValue[uint](1)).GetByValue(total > capacities)
	my.capacities = uint(capacities)

	if my.batches == 0 {
		return nil, ErrBatchInvalid
	}
	if my.capacities == 0 {
		return nil, ErrCapacityInvalid
	}

	var (
		wg      = sync.WaitGroup{}
		results = anySlices.New(anySlices.Cap[Result[T]](total))
	)

	for batch := range my.batches {
		for capacity := range my.capacities {
			wg.Add(1)
			go func(b, c uint) {
				defer wg.Done()
				results.Lock()
				results.Append(my.executeWithRetry(func() Result[T] { return fn(b, c) }))
				results.Unlock()
			}(batch, capacity)
		}
		wg.Wait()
	}

	return results, nil
}

func (my *ResultImpl[T]) GetError() error { return my.Error }

func (my *ResultImpl[T]) GetValue() T { return my.Value }

func (my *ResultImpl[T]) IsOK() bool { return my.Error == nil && !my.Skip }

func (my *ResultImpl[T]) IsSkip() bool { return my.Skip }

func (my *ResultImpl[T]) SetOK(ok bool) { my.Error = nil }

func (my *ResultImpl[T]) SetSkip(skip bool) { my.Skip = skip }
