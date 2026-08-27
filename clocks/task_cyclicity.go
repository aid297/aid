package clocks

import (
	_context "context"
	_errors "errors"
	"fmt"
	_time "time"

	_uuid "github.com/google/uuid"
)

var (
	_             Tasker = (*TaskCyclicityImpl)(nil)
	TaskCyclicity TaskCyclicityImpl
)

type (
	TaskCyclicityImpl struct {
		uuid        _uuid.UUID
		name        string
		interval    _time.Duration
		timeout     _time.Duration
		fn          func(tasker Tasker)
		closeCh     chan es
		immediately bool
	}
)

func (*TaskCyclicityImpl) New(interval _time.Duration) *TaskCyclicityImpl {
	return &TaskCyclicityImpl{uuid: _uuid.Must(_uuid.NewV7()), interval: interval, timeout: defaultTimeout, closeCh: make(chan es, 1)}
}

func (*TaskCyclicityImpl) Secondly(seconds uint64) *TaskCyclicityImpl {
	return &TaskCyclicityImpl{uuid: _uuid.Must(_uuid.NewV7()), interval: _time.Duration(seconds) * _time.Second, timeout: defaultTimeout, closeCh: make(chan es, 1)}
}

func (*TaskCyclicityImpl) Minutely(minutes uint64) *TaskCyclicityImpl {
	return &TaskCyclicityImpl{uuid: _uuid.Must(_uuid.NewV7()), interval: _time.Duration(minutes) * _time.Minute, timeout: defaultTimeout, closeCh: make(chan es, 1)}
}

func (*TaskCyclicityImpl) Hourly(hours uint64) *TaskCyclicityImpl {
	return &TaskCyclicityImpl{uuid: _uuid.Must(_uuid.NewV7()), interval: _time.Duration(hours) * _time.Hour, timeout: defaultTimeout, closeCh: make(chan es, 1)}
}

func (*TaskCyclicityImpl) Daily(days uint64) *TaskCyclicityImpl {
	return &TaskCyclicityImpl{uuid: _uuid.Must(_uuid.NewV7()), interval: _time.Duration(days) * _time.Hour * 24, timeout: defaultTimeout, closeCh: make(chan es, 1)}
}

func (*TaskCyclicityImpl) Weekly(weeks uint64) *TaskCyclicityImpl {
	return &TaskCyclicityImpl{uuid: _uuid.Must(_uuid.NewV7()), interval: _time.Duration(weeks) * _time.Hour * 24 * 7, timeout: defaultTimeout, closeCh: make(chan es, 1)}
}

func (my *TaskCyclicityImpl) String() string {
	return fmt.Sprintf("uuid: %s, name: %s, interval: %s, timeout: %s", my.uuid, my.name, my.interval, my.timeout)
}

func (my *TaskCyclicityImpl) UUID() _uuid.UUID { return my.uuid }

func (my *TaskCyclicityImpl) SetName(name string) Tasker { my.name = name; return my }

func (my *TaskCyclicityImpl) Name() string { return my.name }

func (my *TaskCyclicityImpl) SetTimeout(timeout _time.Duration) Tasker {
	if timeout > 0 {
		my.timeout = timeout
	}

	return my
}

func (my *TaskCyclicityImpl) Timeout() _time.Duration { return my.timeout }

func (my *TaskCyclicityImpl) SetFn(fn func(tasker Tasker)) Tasker { my.fn = fn; return my }

func (my *TaskCyclicityImpl) Fn() func(tasker Tasker) { return my.fn }

func (my *TaskCyclicityImpl) SetImmediately(immediately bool) Tasker {
	my.immediately = immediately
	return my
}

func (my *TaskCyclicityImpl) Immediately() bool { return my.immediately }

func (my *TaskCyclicityImpl) Do() {
	ctx, cancel := _context.WithTimeout(_context.Background(), my.timeout)
	defer cancel()

	done := make(chan struct{})
	go func() { defer close(done); my.fn(my) }()

	select {
	case <-ctx.Done():
		clockIns.errHandler(my, _errors.New("任务执行超时"))
	case <-done:
	}
}

func (my *TaskCyclicityImpl) Begin() error {
	if my.fn == nil {
		return _errors.New("回调方法为空")
	}

	if my.Immediately() {
		my.Do()
	}

	ticker := _time.NewTicker(my.interval)
	defer ticker.Stop()

	for {
		select {
		case <-my.closeCh:
			return nil
		case <-ticker.C:
			my.Do()
		}
	}
}

func (my *TaskCyclicityImpl) Stop() Tasker { my.closeCh <- es{}; return my }
