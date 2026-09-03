package clocks

import (
	_context "context"
	"sync"
	_time "time"

	_uuid "github.com/google/uuid"

	_anyMaps "github.com/aid297/aid/v2/anyMaps"
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
		Boot(ctx _context.Context) *ClockImpl
		Clean() *ClockImpl
	}

	Tasker interface {
		String() string
		UUID() _uuid.UUID
		SetName(name string) Tasker
		Name() string
		SetTimeout(timeout _time.Duration) Tasker
		Timeout() _time.Duration
		SetHandler(fn TaskHandler) Tasker
		Handler() TaskHandler
		SetImmediately(immediately bool) Tasker
		EnableImmediately() Tasker
		DisableImmediately() Tasker
		Immediately() bool
		Begin() error
		Stop() Tasker
	}

	TaskHandler func(tasker Tasker)
	ErrHandler  func(tasker Tasker, err error)

	ClockImpl struct {
		taskers    map[_uuid.UUID]Tasker
		errHandler func(tasker Tasker, err error)
	}
)

func init() { OnceClock() }

func OnceClock() *ClockImpl {
	clockOnce.Do(func() {
		clockIns = &ClockImpl{
			taskers:    make(map[_uuid.UUID]Tasker),
			errHandler: func(Tasker, error) {},
		}
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
	clockLock.RLock()
	defer clockLock.RUnlock()

	return clockIns.taskers[uuid]
}

func (my *ClockImpl) Taskers() []Tasker {
	clockLock.RLock()
	defer clockLock.RUnlock()

	return _anyMaps.New(_anyMaps.Map(my.taskers)).GetValues().ToSlice()
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

func (*ClockImpl) Boot(ctx _context.Context) *ClockImpl {
	clockLock.RLock()
	defer clockLock.RUnlock()

	// 本次 Boot 启动的任务使用局部 WaitGroup 跟踪，关停时等待它们的在跑 handler 收尾；
	// 不能用全局共享的 WaitGroup：一次 Boot 的 Wait 会与后续注册动作的 Add 并发，违反 WaitGroup 的使用约束
	var wg sync.WaitGroup
	for _, tasker := range clockIns.taskers {
		t := tasker
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := t.Begin(); err != nil {
				clockIns.errHandler(t, err)
			}
		}()
	}

	go func() {
		// 等待上下文取消
		<-ctx.Done()

		// 注销全部任务——未触发/等待中的任务收到停止信号后立即退出，不再进入下一轮
		// 注意：这里会停止所有已注册任务（含其他途径注册的），但 Boot 只等待自己启动的任务收尾
		clockIns.Clean()

		// 等待在跑任务收尾——正在执行的 handler 执行完（或达到 Timeout）后 Begin 返回，任务 goroutine 退出
		wg.Wait()
	}()

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
