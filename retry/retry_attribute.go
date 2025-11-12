package retry

import (
	"context"
	"time"
)

type (
	Attributer interface {
		Register(retry *Retry)
	}

	AttrSleep   struct{ sleep time.Duration }
	AttrFn      struct{ fn func() error }
	AttrContext struct{ ctx context.Context }
)

func Sleep(sleep time.Duration) AttrSleep { return AttrSleep{sleep: sleep} }

func (my AttrSleep) Register(retry *Retry) { retry.sleep = my.sleep }

func Fn(fn func() error) AttrFn { return AttrFn{fn: fn} }

func (my AttrFn) Register(retry *Retry) { retry.fn = my.fn }

func Context(ctx context.Context) AttrContext { return AttrContext{ctx: ctx} }

func (my AttrContext) Register(retry *Retry) { retry.ctx = my.ctx }
