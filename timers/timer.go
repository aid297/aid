package timers

import (
	"sync"
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
		GetTask(taskerUUID string) Tasker
		MustGetTask(taskerUUID string) Tasker
		runTask(tasker Tasker)
		RemoveTask(uuid string) Timer
		MustRemoveTask(uuid string) Timer
		Start() Timer
		Stop(uuids ...string) Timer
	}

	TimerImpl struct {
		lock    sync.RWMutex
		timers  map[string]Tasker
		running bool
	}
)

func (*TimerImpl) Once() Timer {
	timerOnce.Do(func() { timerIns = &TimerImpl{timers: make(map[string]Tasker)} })
	return timerIns
}

// AddTask 仅登记任务，不会启动；timer 运行中加入的任务也不会自动执行，需改用 AddTaskAndStart
func (my *TimerImpl) AddTask(tasker Tasker) Timer {
	if tasker == nil {
		return my
	}

	my.lock.Lock()
	defer my.lock.Unlock()

	my.timers[tasker.GetUUID().String()] = tasker
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

func (my *TimerImpl) GetTask(taskerUUID string) Tasker {
	my.lock.RLock()
	defer my.lock.RUnlock()

	return my.timers[taskerUUID]
}

func (my *TimerImpl) MustGetTask(taskerUUID string) Tasker {
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

// RemoveTask 删除任务
func (my *TimerImpl) RemoveTask(uuid string) Timer {
	my.lock.Lock()
	defer my.lock.Unlock()

	if tasker, ok := my.timers[uuid]; ok {
		tasker.Stop()
		delete(my.timers, uuid)
	}

	return my
}

// MustRemoveTask 删除任务
func (my *TimerImpl) MustRemoveTask(uuid string) Timer {
	my.lock.Lock()
	defer my.lock.Unlock()

	if tasker, ok := my.timers[uuid]; ok {
		tasker.Stop()
		delete(my.timers, uuid)
	} else {
		panic("任务不存在")
	}

	return my
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

// Stop 停止全部或指定的任务
func (my *TimerImpl) Stop(uuids ...string) Timer {
	my.lock.Lock()
	my.running = false
	// 锁内快照任务列表，释放锁后再停止任务，避免任务 errHandler 内再调用 AddTask/GetTask 时死锁
	var taskers = make([]Tasker, 0, len(my.timers))

	if len(uuids) > 0 {
		for idx := range uuids {
			if tasker, ok := my.timers[uuids[idx]]; ok {
				taskers = append(taskers, tasker)
			}
		}
	} else {
		for _, tasker := range my.timers {
			taskers = append(taskers, tasker)
		}
	}

	my.lock.Unlock()

	for _, tasker := range taskers {
		tasker.Stop()
	}

	return my
}
