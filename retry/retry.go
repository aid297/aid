package retry

import (
	"context"
	"math/rand"
	"time"
)

type Retry struct {
	sleep time.Duration
	fn    func() error
	ctx   context.Context
}

func (Retry) New(attrs ...Attributer) Retry {
	ins := Retry{fn: nil, ctx: context.TODO()}
	return ins.Set(attrs...)
}

func (my Retry) Set(attrs ...Attributer) Retry {
	if len(attrs) > 0 {
		for idx := range attrs {
			attrs[idx].Register(&my)
		}
	}

	return my
}

// Simple 线性重试
func (my Retry) Simple(attempts int) error {
	if my.fn == nil {
		return nil
	}

	if err := my.fn(); err != nil {
		if attempts--; attempts > 0 {
			time.Sleep(my.sleep)
			return my.Simple(attempts)
		}
		return err
	}

	return nil
}

// Do 指数退避
func (my Retry) Do(attempts int) error {
	if my.fn == nil {
		return nil
	}

	if err := my.fn(); err != nil {
		if attempts--; attempts > 0 {
			time.Sleep(my.sleep)
			return my.Set(Sleep(2 * my.sleep)).Do(attempts)
		}
		return err
	}

	return nil
}

// WithContext 带上下文的重试
func (my Retry) WithContext(attempts int) error {
	if my.fn == nil {
		return nil
	}

	if err := my.fn(); err != nil {
		if attempts--; attempts > 0 {
			select {
			case <-time.After(my.sleep):
				return my.Set(Sleep(2 * my.sleep)).WithContext(attempts) // 指数退避
			case <-my.ctx.Done():
				return my.ctx.Err()
			}
		}
		return err
	}

	return nil
}

func (my Retry) WithContextAndJitter(attempts int) error {
	if my.fn == nil {
		return nil
	}

	if err := my.fn(); err != nil {
		if attempts--; attempts > 0 {
			// 加入随机退避
			jitter := time.Duration(rand.Int63n(int64(my.sleep)))
			my.sleep = my.sleep + jitter

			select {
			case <-time.After(my.sleep):
				return my.Set(Sleep(2 * my.sleep)).WithContextAndJitter(attempts) // 指数退避
			case <-my.ctx.Done():
				return my.ctx.Err()
			}
		}
		return err
	}

	return nil
}
