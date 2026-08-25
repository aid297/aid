package timers

import (
	"errors"
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
		New(fn FN) (TimerTasker, error)
		SetTimeout(timeout time.Duration) TimerTasker
		SetErrHandler(errHandler ErrHandler) TimerTasker
		SetName(name string) TimerTasker
		SetAt(at time.Time) TimerTasker
		SetDelay(delay time.Duration) TimerTasker
		SetDelayAt(at time.Time, loc *time.Location) TimerTasker
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
		runSchedule(interval int) (err error)
		computeNext(loc *time.Location) time.Time
		advanceNext(current time.Time, loc *time.Location, interval int) time.Time
		safeFn() (err error)
		Stop() (err error)
	}

	TimerTaskerImpl struct {
		UUID       _uuid.UUID
		Name       string
		TaskType   TaskTypeTag
		Fn         FN
		Interval   time.Duration
		At         time.Time     // TaskTypeOnce: 到点执行一次的目标时间；周期性任务: 锚点时间，定义每个周期内的精确执行时刻
		Delay      time.Duration // 周期任务启动后的首次延迟时长，在进入第一个调度周期前等待
		stopCh     chan struct{}
		stopOnce   sync.Once
		Timeout    time.Duration
		ErrHandler ErrHandler
	}

	FN func()

	ErrHandler func(tasker TimerTasker, err error)
)

func (*TimerTaskerImpl) New(fn FN) (TimerTasker, error) {
	timerTasker := &TimerTaskerImpl{UUID: _uuid.Must(_uuid.NewV7()), At: defaultAt, Timeout: defaultTimeout, stopCh: make(chan struct{}, 1)}

	if fn == nil {
		return nil, errors.New("回调方法不能为空")
	}

	return timerTasker, nil
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

func (my *TimerTaskerImpl) SetAt(at time.Time) TimerTasker { my.At = at; return my }

func (my *TimerTaskerImpl) SetDelay(delay time.Duration) TimerTasker {
	if delay > 0 {
		my.Delay = delay
	}
	return my
}

// SetDelayAt 根据目标时间和时区计算延迟时长并设置；
// 将 at 和当前时间统一转换到 loc 时区后再计算差值，避免因时区不一致导致计算错误；
// 若计算结果 <= 0（目标时间已过），则不设置 Delay
func (my *TimerTaskerImpl) SetDelayAt(at time.Time, loc *time.Location) TimerTasker {
	if loc == nil {
		loc = time.Local
	}
	var now = time.Now().In(loc)
	var target = at.In(loc)
	var delay = target.Sub(now)
	if delay > 0 {
		my.Delay = delay
	}
	return my
}

func (my *TimerTaskerImpl) Once(at time.Time) TimerTasker {
	my.UUID = _uuid.Must(_uuid.NewV7())
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
			delay = 0 // 目标时间已过，立即执行
		}
		return my.runAfter(delay)
	case TaskTypeDaily, TaskTypeHourly, TaskTypeMinutely:
		return my.runSchedule(int(my.Interval))
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

// runSchedule 通用周期调度：若有 Delay 则先等待，之后计算下一个精确执行时刻，循环等待并执行
func (my *TimerTaskerImpl) runSchedule(interval int) (err error) {
	// 首次延迟：在进入周期调度前等待 Delay 时长，期间可被 Stop 取消
	if my.Delay > 0 {
		var delayTimer = time.NewTimer(my.Delay)
		select {
		case <-delayTimer.C:
		case <-my.stopCh:
			delayTimer.Stop()
			return
		}
		delayTimer.Stop()
	}

	var loc = my.At.Location()
	var next = my.computeNext(loc)

	var timeoutCh <-chan time.Time
	if my.Timeout > 0 {
		timeoutCh = time.After(my.Timeout)
	}

	for {
		var delay = time.Until(next)
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
		next = my.advanceNext(next, loc, interval)
	}
}

// computeNext 根据任务类型和锚点，计算第一个未来执行时刻
func (my *TimerTaskerImpl) computeNext(loc *time.Location) time.Time {
	var now = time.Now().In(loc)
	var h, m, s = my.At.Clock()

	switch my.TaskType {
	// case TaskTypeYearly:
	// 	var next = time.Date(now.Year(), my.At.Month(), my.At.Day(), h, m, s, 0, loc)
	// 	if !next.After(now) {
	// 		next = time.Date(now.Year()+1, my.At.Month(), my.At.Day(), h, m, s, 0, loc)
	// 	}
	// 	return next
	// case TaskTypeMonthly:
	// 	var next = time.Date(now.Year(), now.Month(), my.At.Day(), h, m, s, 0, loc)
	// 	if !next.After(now) {
	// 		next = time.Date(now.Year(), now.Month()+1, my.At.Day(), h, m, s, 0, loc)
	// 	}
	// 	return next
	// case TaskTypeWeekly:
	// 	var next = time.Date(now.Year(), now.Month(), now.Day(), h, m, s, 0, loc)
	// 	var targetWeekday = my.At.Weekday()
	// 	var daysAhead = (int(targetWeekday) - int(now.Weekday()) + 7) % 7
	// 	if daysAhead == 0 && !next.After(now) {
	// 		daysAhead = 7
	// 	}
	// 	next = next.AddDate(0, 0, daysAhead)
	// 	return next
	case TaskTypeDaily:
		var next = time.Date(now.Year(), now.Month(), now.Day(), h, m, s, 0, loc)
		if !next.After(now) {
			next = time.Date(now.Year(), now.Month(), now.Day()+1, h, m, s, 0, loc)
		}
		return next
	case TaskTypeHourly:
		var next = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), m, s, 0, loc)
		if !next.After(now) {
			next = next.Add(time.Hour)
		}
		return next
	case TaskTypeMinutely:
		var next = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), s, 0, loc)
		if !next.After(now) {
			next = next.Add(time.Minute)
		}
		return next
	default:
		return now
	}
}

// advanceNext 按任务类型将当前执行时刻推进到下一个周期
func (my *TimerTaskerImpl) advanceNext(current time.Time, loc *time.Location, interval int) time.Time {
	var h, m, s = my.At.Clock()

	switch my.TaskType {
	// case TaskTypeYearly:
	// 	return time.Date(current.Year()+interval, my.At.Month(), my.At.Day(), h, m, s, 0, loc)
	// case TaskTypeMonthly:
	// 	return time.Date(current.Year(), current.Month()+time.Month(interval), my.At.Day(), h, m, s, 0, loc)
	// case TaskTypeWeekly:
	// 	return current.AddDate(0, 0, 7*interval)
	case TaskTypeDaily:
		return time.Date(current.Year(), current.Month(), current.Day()+interval, h, m, s, 0, loc)
	case TaskTypeHourly:
		return current.Add(time.Duration(interval) * time.Hour)
	case TaskTypeMinutely:
		return current.Add(time.Duration(interval) * time.Minute)
	default:
		return current
	}
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
