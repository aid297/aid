package retries

import (
	_context "context"
	_rand "math/rand"
	_time "time"
)

type (
	Retry interface {
		Set(attrs ...RetryAttr) Retry
		Linear(attempts uint) error
		Exponent(attempts uint) error
		LinearWithContext(attempts uint) error
		JitterWithContext(attempts uint) error
	}

	RetryFn func(retried uint) (err error)

	RetryImpl struct {
		sleep _time.Duration
		fn    RetryFn
		ctx   _context.Context
	}
)

func NewRetry(attrs ...RetryAttr) Retry {
	return (&RetryImpl{fn: nil, ctx: _context.TODO()}).Set(attrs...)
}

// Set 设置属性
func (my *RetryImpl) Set(attrs ...RetryAttr) Retry {
	if len(attrs) > 0 {
		for idx := range attrs {
			attrs[idx](my)
		}
	}

	return my
}

// Linear 线性重试
func (my *RetryImpl) Linear(attempts uint) error {
	if my.fn == nil {
		return nil
	}

	if err := my.fn(attempts); err != nil {
		if attempts--; attempts > 0 {
			_time.Sleep(my.sleep)
			return my.Linear(attempts)
		}
		return err
	}

	return nil
}

// Exponent 指数退避
func (my *RetryImpl) Exponent(attempts uint) error {
	if my.fn == nil {
		return nil
	}

	if err := my.fn(attempts); err != nil {
		if attempts--; attempts > 0 {
			_time.Sleep(my.sleep)
			return my.Set(Sleep(2 * my.sleep)).Exponent(attempts)
		}
		return err
	}

	return nil
}

// LinearWithContext 线性重试携带上下文
func (my *RetryImpl) LinearWithContext(attempts uint) error {
	if my.fn == nil {
		return nil
	}

	if err := my.fn(attempts); err != nil {
		if attempts--; attempts > 0 {
			select {
			case <-_time.After(my.sleep):
				return my.Set(Sleep(2 * my.sleep)).LinearWithContext(attempts) // 指数退避
			case <-my.ctx.Done():
				return my.ctx.Err()
			}
		}
		return err
	}

	return nil
}

// JitterWithContext 随机数退避携带上下文
func (my *RetryImpl) JitterWithContext(attempts uint) error {
	if my.fn == nil {
		return nil
	}

	if err := my.fn(attempts); err != nil {
		if attempts--; attempts > 0 {
			// 加入随机退避
			jitter := _time.Duration(_rand.Int63n(int64(my.sleep)))
			my.sleep = my.sleep + jitter

			select {
			case <-_time.After(my.sleep):
				return my.Set(Sleep(2 * my.sleep)).JitterWithContext(attempts) // 指数退避
			case <-my.ctx.Done():
				return my.ctx.Err()
			}
		}
		return err
	}

	return nil
}
