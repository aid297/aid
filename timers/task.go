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
	TaskTypeCyclicity
	TaskTypeOnce
)

var (
	TaskIntervalDaily  = time.Duration(24) * time.Hour
	TaskIntervalWeekly = time.Duration(7) * TaskIntervalDaily
)

type (
	TaskTypeTag int

	Tasker interface {
		Once(at time.Time, taskerOption TaskerOption) (tasker Tasker, err error)
		After(interval time.Duration, taskerOption TaskerOption) (tasker Tasker, err error)
		Cyclicity(interval time.Duration, taskerOption TaskerOption, at time.Time) (tasker Tasker, err error)
		GetUUID() _uuid.UUID
		GetName() string
		GetTaskType() TaskTypeTag
		GetInterval() time.Duration
		GetTimeout() time.Duration
		GetAt() time.Time
		GetFn() FN
		GetErrHandler() ErrHandler
		Start() (err error)
		Stop() (err error)
	}

	TaskerImpl struct {
		UUID       _uuid.UUID
		Name       string
		TaskType   TaskTypeTag
		Fn         FN
		Interval   time.Duration
		At         time.Time // TaskTypeOnce: 到点执行一次的目标时间；TaskTypeCyclicity: 锚点时间，提取 Location 和 Hour/Min/Sec 实现每日精确调度
		stop       chan struct{}
		stopOnce   sync.Once
		Timeout    time.Duration
		ErrHandler ErrHandler
	}

	TaskerOption struct {
		err        error
		Name       string
		Fn         FN
		Timeout    time.Duration
		ErrHandler ErrHandler
	}

	FN func()

	ErrHandler func(tasker Tasker, err error)
)

func NewTaskerOption(name string, timeout time.Duration, fn FN, errHandler ErrHandler) TaskerOption {
	if fn == nil {
		return TaskerOption{err: errors.New("回调方法不能为空")}
	}
	if errHandler == nil {
		return TaskerOption{err: errors.New("错误处理方法不能为空")}
	}

	return TaskerOption{Name: name, Timeout: timeout, Fn: fn, ErrHandler: errHandler, err: nil}
}

func (*TaskerImpl) Once(at time.Time, taskerOption TaskerOption) (tasker Tasker, err error) {
	if taskerOption.err != nil {
		return nil, taskerOption.err
	}

	if at.IsZero() {
		return nil, errors.New("启动时间错误")
	}

	tasker = &TaskerImpl{
		UUID:       _uuid.Must(_uuid.NewV7()),
		TaskType:   TaskTypeOnce,
		Name:       taskerOption.Name,
		At:         at,
		Timeout:    taskerOption.Timeout,
		Fn:         taskerOption.Fn,
		ErrHandler: taskerOption.ErrHandler,
		stop:       make(chan struct{}, 1),
	}

	return
}

func (*TaskerImpl) After(interval time.Duration, taskerOption TaskerOption) (tasker Tasker, err error) {
	if taskerOption.err != nil {
		return nil, taskerOption.err
	}

	if interval <= 0 {
		return nil, errors.New("间隔时间错误")
	}

	tasker = &TaskerImpl{
		UUID:       _uuid.Must(_uuid.NewV7()),
		TaskType:   TaskTypeAfter,
		Name:       taskerOption.Name,
		Interval:   interval,
		Timeout:    taskerOption.Timeout,
		Fn:         taskerOption.Fn,
		ErrHandler: taskerOption.ErrHandler,
		stop:       make(chan struct{}, 1),
	}

	return
}

func (*TaskerImpl) Cyclicity(interval time.Duration, taskerOption TaskerOption, at time.Time) (tasker Tasker, err error) {
	if taskerOption.err != nil {
		return nil, taskerOption.err
	}

	if interval <= 0 {
		return nil, errors.New("间隔时间错误")
	}

	tasker = &TaskerImpl{
		UUID:       _uuid.Must(_uuid.NewV7()),
		TaskType:   TaskTypeCyclicity,
		Name:       taskerOption.Name,
		Interval:   interval,
		Timeout:    taskerOption.Timeout,
		Fn:         taskerOption.Fn,
		ErrHandler: taskerOption.ErrHandler,
		stop:       make(chan struct{}, 1),
		At:         at,
	}

	return
}

func (my *TaskerImpl) GetUUID() _uuid.UUID { return my.UUID }

func (my *TaskerImpl) GetName() string { return my.Name }

func (my *TaskerImpl) GetTaskType() TaskTypeTag { return my.TaskType }

func (my *TaskerImpl) GetInterval() time.Duration { return my.Interval }

func (my *TaskerImpl) GetTimeout() time.Duration { return my.Timeout }

func (my *TaskerImpl) GetAt() time.Time { return my.At }

func (my *TaskerImpl) GetFn() FN { return my.Fn }

func (my *TaskerImpl) GetErrHandler() ErrHandler { return my.ErrHandler }

func (my *TaskerImpl) Start() (err error) {
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
	case TaskTypeCyclicity:
		// At 非零时启用锚点调度：从 At 提取时区与时:分:秒，每天在精确时刻执行
		if !my.At.IsZero() {
			return my.runCyclicityAt()
		}

		var ticker = time.NewTicker(my.Interval)
		defer ticker.Stop()

		var timeoutCh <-chan time.Time
		if my.Timeout > 0 {
			timeoutCh = time.After(my.Timeout) // timeout <= 0 时 timeoutCh 为 nil，该 case 永不触发
		}

		for {
			select {
			case <-ticker.C: // 周期触发
				select {
				case <-my.stop: // 已停止则不再执行
					return
				default:
				}
				// 单次执行的 panic 不中断循环调度
				if fnErr := my.safeFn(); fnErr != nil {
					my.ErrHandler(my, fnErr)
				}
			case <-my.stop: // 手动停止，不算错误
				return
			case <-timeoutCh: // 等待超时
				my.ErrHandler(my, fmt.Errorf("定时任务执行超时：%s", my.Timeout))
				return
			}
		}
	default:
		return fmt.Errorf("不支持的定时任务类型 %d", my.TaskType)
	}
}

// runAfter 在 delay 之后执行一次回调，供 After/Once 共用
func (my *TaskerImpl) runAfter(delay time.Duration) (err error) {
	var c = make(chan error, 1)
	var timer = time.AfterFunc(delay, func() {
		select {
		case <-my.stop: // 已停止则不再执行
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
	case <-my.stop: // 手动停止，不算错误
		timer.Stop()
		return
	case <-timeoutCh: // 等待超时
		timer.Stop()
		my.ErrHandler(my, fmt.Errorf("定时任务执行超时：%s", my.Timeout))
		return
	}
	return
}

// runCyclicityAt 锚点调度：从 At 提取时区与时:分:秒，每隔 Interval 在精确时刻执行一次；
// 每次触发都在目标时区重新构造日期，确保 DST 切换等场景下墙钟时间始终正确
func (my *TaskerImpl) runCyclicityAt() (err error) {
	var loc = my.At.Location()
	var h, m, s = my.At.Clock()
	var days = int(my.Interval / TaskIntervalDaily)
	if days < 1 {
		days = 1
	}

	// 从今天（锚点日期）的 H:M:S 开始，若已过则推进到下一个周期
	var next = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), h, m, s, 0, loc)
	if next.Before(time.Now()) {
		next = time.Date(my.At.Year(), my.At.Month(), my.At.Day(), h, m, s, 0, loc)
		for !next.After(time.Now()) {
			next = next.AddDate(0, 0, days)
		}
	}

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
			case <-my.stop:
				timer.Stop()
				return
			default:
			}
			if fnErr := my.safeFn(); fnErr != nil {
				my.ErrHandler(my, fnErr)
			}
		case <-my.stop:
			timer.Stop()
			return
		case <-timeoutCh:
			timer.Stop()
			my.ErrHandler(my, fmt.Errorf("定时任务执行超时：%s", my.Timeout))
			return
		}

		timer.Stop()
		// 在目标时区重新构造日期，保证跨 DST 边界时墙钟时间不变
		next = next.AddDate(0, 0, days)
		next = time.Date(next.Year(), next.Month(), next.Day(), h, m, s, 0, loc)
	}
}

// safeFn 带 panic 保护地执行回调，供循环任务使用，单次 panic 不中断调度
func (my *TaskerImpl) safeFn() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("定时任务执行失败：%v", r)
		}
	}()

	my.Fn()
	return
}

func (my *TaskerImpl) Stop() (err error) { my.stopOnce.Do(func() { close(my.stop) }); return }
