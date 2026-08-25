package timers

import (
	"fmt"
	"sync"
	"time"

	_uuid "github.com/google/uuid"
)

const (
	TaskTypeAfter TaskTypeTag = iota
	TaskTypeDaily
	TaskTypeHourly
	TaskTypeMinutely
	TaskTypeOnce
)

var (
	defaultTimeout = 24 * time.Hour
	defaultAfter   = 30 * time.Second
	defaultAt      = time.Date(0, 0, 0, 0, 0, 0, 0, time.Local)
)

type (
	TaskTypeTag int

	TimerTasker interface {
		New(fn FN) TimerTasker
		SetTimeout(timeout time.Duration) TimerTasker
		SetErrHandler(errHandler ErrHandler) TimerTasker
		SetName(name string) TimerTasker
		SetDelay(delay time.Duration) TimerTasker
		SetDelayAt(at time.Time) TimerTasker
		Once(at time.Time) TimerTasker
		After(interval time.Duration) TimerTasker
		Daily(interval int) TimerTasker
		Hourly(interval int) TimerTasker
		Minutely(interval int) TimerTasker
		GetUUID() _uuid.UUID
		GetName() string
		GetTaskType() TaskTypeTag
		GetInterval() time.Duration
		GetTimeout() time.Duration
		GetAt() time.Time
		GetDelay() time.Duration
		GetFn() FN
		GetErrHandler() ErrHandler
		Start() (err error)
		runAfter(delay time.Duration) (err error)
		runSchedule() (err error)
		advanceNext(current time.Time) time.Time
		safeFn() (err error)
		Stop() (err error)
	}

	TimerTaskerImpl struct {
		UUID       _uuid.UUID
		Name       string
		TaskType   TaskTypeTag
		Fn         FN
		Interval   time.Duration
		At         time.Time     // 首次执行的目标时间；Once: 到点执行一次；周期任务: 首次执行时刻
		Delay      time.Duration // 首次延迟时长（SetDelay 设置，保留字段）
		stopCh     chan struct{}
		stopOnce   sync.Once
		Timeout    time.Duration
		ErrHandler ErrHandler
	}

	FN func()

	ErrHandler func(tasker TimerTasker, err error)
)

func (*TimerTaskerImpl) New(fn FN) TimerTasker {
	timerTasker := &TimerTaskerImpl{UUID: _uuid.Must(_uuid.NewV7()), At: defaultAt, Timeout: defaultTimeout, stopCh: make(chan struct{}, 1)}

	if fn == nil {
		timerTasker.Fn = func() {}
	} else {
		timerTasker.Fn = fn
	}

	return timerTasker
}

func (my *TimerTaskerImpl) SetTimeout(timeout time.Duration) TimerTasker {
	if timeout <= 0 {
		my.Timeout = defaultTimeout
	} else {
		my.Timeout = timeout
	}
	return my
}

func (my *TimerTaskerImpl) SetErrHandler(errHandler ErrHandler) TimerTasker {
	my.ErrHandler = errHandler
	return my
}

func (my *TimerTaskerImpl) SetName(name string) TimerTasker { my.Name = name; return my }

// SetDelay 设置单次任务的延迟时长，等同 After
func (my *TimerTaskerImpl) SetDelay(delay time.Duration) TimerTasker {
	return my.After(delay)
}

// SetDelayAt 设置周期任务的首次执行时间；At = at，之后按 Interval 循环
func (my *TimerTaskerImpl) SetDelayAt(at time.Time) TimerTasker {
	my.At = at
	return my
}

func (my *TimerTaskerImpl) Once(at time.Time) TimerTasker {
	my.TaskType = TaskTypeOnce
	my.At = at

	return my
}

func (my *TimerTaskerImpl) After(interval time.Duration) TimerTasker {
	if interval <= 0 {
		my.Interval = defaultAfter
	} else {
		my.Interval = interval
	}

	return my
}

func (my *TimerTaskerImpl) Daily(interval int) TimerTasker {
	if interval < 1 {
		interval = 1
	}

	my.TaskType = TaskTypeDaily
	my.Interval = time.Duration(interval) * (time.Hour * 24)

	return my
}

func (my *TimerTaskerImpl) Hourly(interval int) TimerTasker {
	if interval < 1 {
		interval = 1
	}

	my.TaskType = TaskTypeHourly
	my.Interval = time.Duration(interval) * time.Hour

	return my
}

func (my *TimerTaskerImpl) Minutely(interval int) TimerTasker {
	if interval < 1 {
		interval = 1
	}

	my.TaskType = TaskTypeDaily
	my.Interval = time.Duration(interval) * time.Minute

	return my
}

func (my *TimerTaskerImpl) GetUUID() _uuid.UUID { return my.UUID }

func (my *TimerTaskerImpl) GetName() string { return my.Name }

func (my *TimerTaskerImpl) GetTaskType() TaskTypeTag { return my.TaskType }

func (my *TimerTaskerImpl) GetInterval() time.Duration { return my.Interval }

func (my *TimerTaskerImpl) GetTimeout() time.Duration { return my.Timeout }

func (my *TimerTaskerImpl) GetAt() time.Time { return my.At }

func (my *TimerTaskerImpl) GetDelay() time.Duration { return my.Delay }

func (my *TimerTaskerImpl) GetFn() FN { return my.Fn }

func (my *TimerTaskerImpl) GetErrHandler() ErrHandler { return my.ErrHandler }

func (my *TimerTaskerImpl) Start() (err error) {
	if my.Fn == nil {
		return fmt.Errorf("定时任务执行失败，没有需要执行的方法")
	}

	switch my.TaskType {
	case TaskTypeAfter:
		return my.runAfter(my.Interval)
	case TaskTypeOnce:
		var delay = time.Until(my.At)
		if delay < 0 {
			delay = 0
		}
		return my.runAfter(delay)
	case TaskTypeDaily, TaskTypeHourly, TaskTypeMinutely:
		return my.runSchedule()
	default:
		return fmt.Errorf("不支持的定时任务类型 %d", my.TaskType)
	}
}

// runAfter 在 delay 之后执行一次回调，供 After/Once 共用
func (my *TimerTaskerImpl) runAfter(delay time.Duration) (err error) {
	var c = make(chan error, 1)
	var timer = time.AfterFunc(delay, func() {
		select {
		case <-my.stopCh: // 已停止则不再执行
			return
		default:
		}
		defer func() {
			if r := recover(); r != nil {
				c <- fmt.Errorf("定时任务执行失败：%v", r)
			}
		}()
		my.Fn()
		c <- nil // 缓冲 channel，不会阻塞泄漏
	})

	var timeoutCh <-chan time.Time
	if my.Timeout > 0 {
		timeoutCh = time.After(my.Timeout) // timeout <= 0 时 timeoutCh 为 nil，该 case 永不触发
	}
	select {
	case err = <-c: // 执行完成，失败时上报 errHandler
		if err != nil {
			my.ErrHandler(my, err)
		}
	case <-my.stopCh: // 手动停止，不算错误
		timer.Stop()
		return
	case <-timeoutCh: // 等待超时
		timer.Stop()
		my.ErrHandler(my, fmt.Errorf("定时任务执行超时：%s", my.Timeout))
		return
	}
	return
}

// runSchedule 通用周期调度：从 At 开始，按 Interval 周期执行
func (my *TimerTaskerImpl) runSchedule() (err error) {
	var next = my.computeNext()

	var timeoutCh <-chan time.Time
	if my.Timeout > 0 {
		timeoutCh = time.After(my.Timeout)
	}

	for {
		var delay = time.Until(next)
		if delay < 0 {
			delay = 0
		}
		var timer = time.NewTimer(delay)

		select {
		case <-timer.C:
			select {
			case <-my.stopCh:
				timer.Stop()
				return
			default:
			}
			if fnErr := my.safeFn(); fnErr != nil {
				my.ErrHandler(my, fnErr)
			}
		case <-my.stopCh:
			timer.Stop()
			return
		case <-timeoutCh:
			timer.Stop()
			my.ErrHandler(my, fmt.Errorf("定时任务执行超时：%s", my.Timeout))
			return
		}

		timer.Stop()
		next = my.advanceNext(next)
	}
}

// computeNext 计算首次执行时刻：At 在未来则用 At，否则立即执行
func (my *TimerTaskerImpl) computeNext() time.Time {
	if my.At.After(time.Now()) {
		return my.At
	}
	return time.Now()
}

// advanceNext 将当前执行时刻推进到下一个周期
func (my *TimerTaskerImpl) advanceNext(current time.Time) time.Time {
	return current.Add(my.Interval)
}

// safeFn 带 panic 保护地执行回调，供循环任务使用，单次 panic 不中断调度
func (my *TimerTaskerImpl) safeFn() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("定时任务执行失败：%v", r)
		}
	}()

	my.Fn()
	return
}

func (my *TimerTaskerImpl) Stop() (err error) { my.stopOnce.Do(func() { close(my.stopCh) }); return }
