package clocks

import (
	"sync"
	_time "time"

	_uuid "github.com/google/uuid"
)

var (
	defaultInterval = 1 * _time.Minute
	defaultTimeout  = 24 * _time.Hour
	clockOnce       = sync.Once{}
	clockIns        ClockImpl
	clockLock       = sync.RWMutex{}

	Clock ClockImpl
)

type (
	es struct{}

	Tasker interface {
		String() string
		UUID() _uuid.UUID
		SetName(name string) Tasker
		Name() string
		SetTimeout(timeout _time.Duration) Tasker
		Timeout() _time.Duration
		SetFn(fn func(tasker Tasker)) Tasker
		Fn() func(tasker Tasker)
		Begin() error
		Stop() Tasker
	}

	ClockImpl struct {
		taskers    map[_uuid.UUID]Tasker
		errHandler func(tasker Tasker, err error)
	}
)

func (ClockImpl) Ins() ClockImpl {
	clockOnce.Do(func() { clockIns = ClockImpl{taskers: make(map[_uuid.UUID]Tasker)} })
	return clockIns
}

func (ClockImpl) SetErrHandler(errHandler func(tasker Tasker, err error)) ClockImpl {
	clockLock.Lock()
	defer clockLock.Unlock()

	clockIns.errHandler = errHandler

	return clockIns
}

func (ClockImpl) AddTasker(taskers ...Tasker) ClockImpl {
	clockLock.Lock()
	defer clockLock.Unlock()

	for idx := range taskers {
		clockIns.taskers[taskers[idx].UUID()] = taskers[idx]
	}

	return clockIns
}

func (ClockImpl) AddTaskerAndBegin(tasker Tasker) ClockImpl {
	clockLock.Lock()
	defer clockLock.Unlock()

	clockIns.taskers[tasker.UUID()] = tasker
	clockIns.begin(tasker.UUID())

	return clockIns
}

func (ClockImpl) Tasker(uuid _uuid.UUID) Tasker {
	clockLock.RUnlock()
	defer clockLock.RUnlock()

	return clockIns.taskers[uuid]
}

func (ClockImpl) begin(uuid _uuid.UUID) {
	tasker, ok := clockIns.taskers[uuid]
	if ok {
		go func() {
			if err := tasker.Begin(); err != nil {
				clockIns.errHandler(tasker, err)
			}
		}()
	}
}

func (ClockImpl) Begin(uuid _uuid.UUID) {
	clockLock.RLock()
	defer clockLock.RUnlock()

	clockIns.begin(uuid)
}

func (ClockImpl) close(uuid _uuid.UUID) {
	tasker, ok := clockIns.taskers[uuid]
	if ok {
		tasker.Stop()
		delete(clockIns.taskers, uuid)
	}
}

func (ClockImpl) Close(uuid _uuid.UUID) {
	clockLock.Lock()
	defer clockLock.Unlock()

	clockIns.close(uuid)
}

func (ClockImpl) Boot() ClockImpl {
	clockLock.RLock()
	defer clockLock.RUnlock()

	for _, tasker := range clockIns.taskers {
		clockIns.begin(tasker.UUID())
	}

	return clockIns
}

func (ClockImpl) Clean() ClockImpl {
	clockLock.Lock()
	defer clockLock.Unlock()

	for uuid := range clockIns.taskers {
		clockIns.close(uuid)
	}

	return clockIns
}
