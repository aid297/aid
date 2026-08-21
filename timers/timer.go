package timers

import (
	"sync"

	_uuid "github.com/google/uuid"
)

var (
	Task Tasker = (*TaskerImpl)(nil)
	Time Timer  = (*TimerImpl)(nil)

	timerOnce sync.Once
	timerIns  *TimerImpl
)

type (
	Timer interface {
		Once() Timer
		AddTask(tasker Tasker) Timer
		AddTaskAndStart(tasker Tasker) Timer
		GetTask(taskerUUID _uuid.UUID) Tasker
		MustGetTask(taskerUUID _uuid.UUID) Tasker
		Start() Timer
		Stop() Timer
	}

	TimerImpl struct {
		lock    sync.RWMutex
		timers  map[_uuid.UUID]Tasker
		running bool
	}
)

func (*TimerImpl) Once() Timer {
	timerOnce.Do(func() { timerIns = &TimerImpl{timers: make(map[_uuid.UUID]Tasker)} })
	return timerIns
}

// AddTask 仅登记任务，不会启动；timer 运行中加入的任务也不会自动执行，需改用 AddTaskAndStart
func (my *TimerImpl) AddTask(tasker Tasker) Timer {
	if tasker == nil {
		return my
	}

	my.lock.Lock()
	defer my.lock.Unlock()

	my.timers[tasker.GetUUID()] = tasker
	return my
}

// AddTaskAndStart 添加并立即启动任务，适用于 timer 运行中动态加任务；
// 注意：在 Start 之前调用时，后续 Start 会再次启动该任务，一次性任务会被执行两次
func (my *TimerImpl) AddTaskAndStart(tasker Tasker) Timer {
	if tasker == nil {
		return my
	}

	go my.runTask(tasker)

	return my.AddTask(tasker)
}

func (my *TimerImpl) GetTask(taskerUUID _uuid.UUID) Tasker {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.timers[taskerUUID]
}

func (my *TimerImpl) MustGetTask(taskerUUID _uuid.UUID) Tasker {
	tasker := my.GetTask(taskerUUID)

	if tasker == nil {
		panic("任务不存在")
	}

	return tasker
}

// runTask 在独立协程中执行任务，并带 panic 保护；任务错误由各任务自身的 errHandler 处理
func (my *TimerImpl) runTask(tasker Tasker) {
	defer func() {
		recover()
	}()

	tasker.Start()
}

// Start 启动当前已登记的全部任务；一次性快照语义，启动后再加入的任务不会被启动，且重复调用直接返回
func (my *TimerImpl) Start() Timer {
	my.lock.Lock()

	if my.running {
		my.lock.Unlock()
		return my
	}

	my.running = true
	// 锁内快照任务列表，释放锁后再派生协程，避免任务 errHandler 内再调用 AddTask/Start 时与锁竞争
	var taskers = make([]Tasker, 0, len(my.timers))
	for _, tasker := range my.timers {
		taskers = append(taskers, tasker)
	}
	my.lock.Unlock()

	for _, tasker := range taskers {
		go my.runTask(tasker)
	}

	return my
}

func (my *TimerImpl) Stop() Timer {
	my.lock.Lock()
	my.running = false
	// 锁内快照任务列表，释放锁后再停止任务，避免任务 errHandler 内再调用 AddTask/GetTask 时死锁
	var taskers = make([]Tasker, 0, len(my.timers))
	for _, tasker := range my.timers {
		taskers = append(taskers, tasker)
	}
	my.lock.Unlock()

	for _, tasker := range taskers {
		tasker.Stop()
	}

	return my
}
