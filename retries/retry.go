package retries

import (
	"context"
	"math/rand"
	"time"
)

type (
	Retry interface {
		New(attrs ...Attributer) Retry
		Set(attrs ...Attributer) Retry
		Linear(attempts uint) error
		Exponent(attempts uint) error
		LinearWithContext(attempts uint) error
		JitterWithContext(attempts uint) error
	}

	RetryFn func(retried uint) (err error)

	RetryImpl struct {
		sleep time.Duration
		fn    RetryFn
		ctx   context.Context
	}
)

func NewRetry(attrs ...Attributer) Retry {
	return (&RetryImpl{fn: nil, ctx: context.TODO()}).Set(attrs...)
}

// New 实例化
func (*RetryImpl) New(attrs ...Attributer) Retry { return NewRetry(attrs...) }

// Set 设置属性
func (my *RetryImpl) Set(attrs ...Attributer) Retry {
	if len(attrs) > 0 {
		for idx := range attrs {
			attrs[idx].Register(my)
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
			time.Sleep(my.sleep)
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
			time.Sleep(my.sleep)
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
			case <-time.After(my.sleep):
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
			jitter := time.Duration(rand.Int63n(int64(my.sleep)))
			my.sleep = my.sleep + jitter

			select {
			case <-time.After(my.sleep):
				return my.Set(Sleep(2 * my.sleep)).JitterWithContext(attempts) // 指数退避
			case <-my.ctx.Done():
				return my.ctx.Err()
			}
		}
		return err
	}

	return nil
}
