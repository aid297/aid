package clocks

import (
	_context "context"
	_errors "errors"
	"fmt"
	_time "time"

	_uuid "github.com/google/uuid"
)

var (
	_        Tasker = (*TaskOnceImpl)(nil)
	TaskOnce TaskOnceImpl
)

type TaskOnceImpl struct {
	uuid     _uuid.UUID
	name     string
	interval _time.Duration
	timeout  _time.Duration
	fn       func(tasker Tasker)
	closeCh  chan es
}

func (TaskOnceImpl) After(interval _time.Duration) TaskOnceImpl {
	if interval <= 0 {
		interval = defaultInterval
	}
	return TaskOnceImpl{uuid: _uuid.Must(_uuid.NewV7()), interval: interval, timeout: defaultTimeout, closeCh: make(chan es, 1)}
}

func (TaskOnceImpl) At(time _time.Time, loc *_time.Location) TaskOnceImpl {
	return TaskOnce.After(time.Sub(_time.Now().In(loc)))
}

func (my TaskOnceImpl) String() string {
	return fmt.Sprintf("uuid: %s, name: %s, interval: %s, timeout: %s", my.uuid, my.name, my.interval, my.timeout)
}

func (my TaskOnceImpl) UUID() _uuid.UUID { return my.uuid }

func (my TaskOnceImpl) SetName(name string) Tasker { my.name = name; return my }

func (my TaskOnceImpl) Name() string { return my.name }

func (my TaskOnceImpl) SetTimeout(timeout _time.Duration) Tasker {
	if timeout > 0 {
		my.timeout = timeout
	}
	return my
}

func (my TaskOnceImpl) Timeout() _time.Duration { return my.timeout }

func (my TaskOnceImpl) SetFn(fn func(tasker Tasker)) Tasker { my.fn = fn; return my }

func (my TaskOnceImpl) Fn() func(tasker Tasker) { return my.fn }

func (my TaskOnceImpl) Begin() error {
	if my.fn == nil {
		return _errors.New("回调方法为空")
	}

	ctx, cancel := _context.WithTimeout(_context.Background(), my.timeout)
	defer cancel()

	timer := _time.NewTimer(my.interval)

	select {
	case <-ctx.Done():
		return _errors.New("任务执行超时")
	case <-my.closeCh:
		return nil
	case <-timer.C:
		my.fn(my)
		return nil
	}
}

func (my TaskOnceImpl) Stop() Tasker { my.closeCh <- es{}; return my }
