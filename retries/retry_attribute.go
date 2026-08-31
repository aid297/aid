package retries

import (
	_context "context"
	_time "time"
)

type RetryAttr func(retry *RetryImpl)

func Sleep(sleep _time.Duration) RetryAttr {
	return func(retry *RetryImpl) { retry.sleep = sleep }
}

func Fn(fn RetryFn) RetryAttr {
	return func(retry *RetryImpl) { retry.fn = fn }
}

func Context(ctx _context.Context) RetryAttr {
	return func(retry *RetryImpl) { retry.ctx = ctx }
}
