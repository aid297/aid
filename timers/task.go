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
		Cyclicity(interval time.Duration, taskerOption TaskerOption) (tasker Tasker, err error)
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
		At         time.Time // 仅 TaskTypeOnce 使用：到点执行一次的目标时间
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

func (*TaskerImpl) Cyclicity(interval time.Duration, taskerOption TaskerOption) (tasker Tasker, err error) {
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
