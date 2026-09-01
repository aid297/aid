package clocks

import (
	"sync"
	_time "time"

	_uuid "github.com/google/uuid"

	"github.com/aid297/aid/v2/anyMaps"
)

var (
	defaultInterval = 1 * _time.Minute
	defaultTimeout  = 24 * _time.Hour
	clockOnce       = sync.Once{}
	clockIns        *ClockImpl
	clockLock       = sync.RWMutex{}

	_ Clock = (*ClockImpl)(nil)
)

type (
	es struct{}

	Clock interface {
		SetErrHandler(errHandler func(tasker Tasker, err error)) Clock
		AddTasker(taskers ...Tasker) Clock
		AddTaskerAndBegin(tasker Tasker) Clock
		Tasker(uuid _uuid.UUID) Tasker
		Taskers() []Tasker
		DeleteTasker(uuids ..._uuid.UUID) Clock
		begin(uuid _uuid.UUID)
		Begin(uuid _uuid.UUID)
		close(uuid _uuid.UUID)
		Close(uuid _uuid.UUID)
		Boot() *ClockImpl
		Clean() *ClockImpl
	}

	Tasker interface {
		String() string
		UUID() _uuid.UUID
		SetName(name string) Tasker
		Name() string
		SetTimeout(timeout _time.Duration) Tasker
		Timeout() _time.Duration
		SetFn(fn func(tasker Tasker)) Tasker
		Fn() func(tasker Tasker)
		SetImmediately(immediately bool) Tasker
		EnableImmediately() Tasker
		DisableImmediately() Tasker
		Immediately() bool
		Begin() error
		Stop() Tasker
	}

	ClockImpl struct {
		taskers    map[_uuid.UUID]Tasker
		errHandler func(tasker Tasker, err error)
	}
)

func init() { OnceClock() }

func OnceClock() *ClockImpl {
	clockOnce.Do(func() {
		clockIns = &ClockImpl{taskers: make(map[_uuid.UUID]Tasker), errHandler: func(Tasker, error) {}}
	})
	return clockIns
}

func (*ClockImpl) SetErrHandler(errHandler func(tasker Tasker, err error)) Clock {
	clockLock.Lock()
	defer clockLock.Unlock()

	clockIns.errHandler = errHandler

	return clockIns
}

func (*ClockImpl) AddTasker(taskers ...Tasker) Clock {
	clockLock.Lock()
	defer clockLock.Unlock()

	for idx := range taskers {
		clockIns.taskers[taskers[idx].UUID()] = taskers[idx]
	}

	return clockIns
}

func (*ClockImpl) AddTaskerAndBegin(tasker Tasker) Clock {
	clockLock.Lock()
	defer clockLock.Unlock()

	clockIns.taskers[tasker.UUID()] = tasker
	clockIns.begin(tasker.UUID())

	return clockIns
}

func (*ClockImpl) Tasker(uuid _uuid.UUID) Tasker {
	clockLock.RUnlock()
	defer clockLock.RUnlock()

	return clockIns.taskers[uuid]
}

func (my *ClockImpl) Taskers() []Tasker {
	clockLock.RLock()
	defer clockLock.RUnlock()

	return anyMaps.New(anyMaps.Map(my.taskers)).GetValues().ToSlice()
}

func (*ClockImpl) DeleteTasker(uuids ..._uuid.UUID) Clock {
	clockLock.Lock()
	defer clockLock.Unlock()

	for idx := range uuids {
		delete(clockIns.taskers, uuids[idx])
	}

	return clockIns
}

func (*ClockImpl) begin(uuid _uuid.UUID) {
	tasker, ok := clockIns.taskers[uuid]
	if ok {
		go func() {
			if err := tasker.Begin(); err != nil {
				clockIns.errHandler(tasker, err)
			}
		}()
	}
}

func (*ClockImpl) Begin(uuid _uuid.UUID) {
	clockLock.RLock()
	defer clockLock.RUnlock()

	clockIns.begin(uuid)
}

func (*ClockImpl) close(uuid _uuid.UUID) {
	tasker, ok := clockIns.taskers[uuid]
	if ok {
		tasker.Stop()
		delete(clockIns.taskers, uuid)
	}
}

func (*ClockImpl) Close(uuid _uuid.UUID) {
	clockLock.Lock()
	defer clockLock.Unlock()

	clockIns.close(uuid)
}

func (*ClockImpl) Boot() *ClockImpl {
	clockLock.RLock()
	defer clockLock.RUnlock()

	for _, tasker := range clockIns.taskers {
		clockIns.begin(tasker.UUID())
	}

	return clockIns
}

func (*ClockImpl) Clean() *ClockImpl {
	clockLock.Lock()
	defer clockLock.Unlock()

	for uuid := range clockIns.taskers {
		clockIns.close(uuid)
	}

	return clockIns
}
