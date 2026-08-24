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
	TaskTypeYearly
	TaskTypeMonthly
	TaskTypeWeekly
	TaskTypeDaily
	TaskTypeHourly
	TaskTypeMinutely
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
		Yearly(interval int, at time.Time, taskerOption TaskerOption) (tasker Tasker, err error)
		Monthly(interval int, at time.Time, taskerOption TaskerOption) (tasker Tasker, err error)
		Weekly(interval int, at time.Time, taskerOption TaskerOption) (tasker Tasker, err error)
		Daily(interval int, at time.Time, taskerOption TaskerOption) (tasker Tasker, err error)
		Hourly(interval int, at time.Time, taskerOption TaskerOption) (tasker Tasker, err error)
		Minutely(interval int, at time.Time, taskerOption TaskerOption) (tasker Tasker, err error)
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
		At         time.Time // TaskTypeOnce: 到点执行一次的目标时间；周期性任务: 锚点时间，定义每个周期内的精确执行时刻
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

func (*TaskerImpl) Yearly(interval int, at time.Time, taskerOption TaskerOption) (tasker Tasker, err error) {
	if taskerOption.err != nil {
		return nil, taskerOption.err
	}
	if at.IsZero() {
		return nil, errors.New("锚点时间不能为空")
	}
	if interval < 1 {
		interval = 1
	}

	tasker = &TaskerImpl{
		UUID:       _uuid.Must(_uuid.NewV7()),
		TaskType:   TaskTypeYearly,
		Name:       taskerOption.Name,
		Interval:   time.Duration(interval),
		At:         at,
		Timeout:    taskerOption.Timeout,
		Fn:         taskerOption.Fn,
		ErrHandler: taskerOption.ErrHandler,
		stop:       make(chan struct{}, 1),
	}
	return
}

func (*TaskerImpl) Monthly(interval int, at time.Time, taskerOption TaskerOption) (tasker Tasker, err error) {
	if taskerOption.err != nil {
		return nil, taskerOption.err
	}
	if at.IsZero() {
		return nil, errors.New("锚点时间不能为空")
	}
	if interval < 1 {
		interval = 1
	}

	tasker = &TaskerImpl{
		UUID:       _uuid.Must(_uuid.NewV7()),
		TaskType:   TaskTypeMonthly,
		Name:       taskerOption.Name,
		Interval:   time.Duration(interval),
		At:         at,
		Timeout:    taskerOption.Timeout,
		Fn:         taskerOption.Fn,
		ErrHandler: taskerOption.ErrHandler,
		stop:       make(chan struct{}, 1),
	}
	return
}

func (*TaskerImpl) Weekly(interval int, at time.Time, taskerOption TaskerOption) (tasker Tasker, err error) {
	if taskerOption.err != nil {
		return nil, taskerOption.err
	}
	if at.IsZero() {
		return nil, errors.New("锚点时间不能为空")
	}
	if interval < 1 {
		interval = 1
	}

	tasker = &TaskerImpl{
		UUID:       _uuid.Must(_uuid.NewV7()),
		TaskType:   TaskTypeWeekly,
		Name:       taskerOption.Name,
		Interval:   time.Duration(interval),
		At:         at,
		Timeout:    taskerOption.Timeout,
		Fn:         taskerOption.Fn,
		ErrHandler: taskerOption.ErrHandler,
		stop:       make(chan struct{}, 1),
	}
	return
}

func (*TaskerImpl) Daily(interval int, at time.Time, taskerOption TaskerOption) (tasker Tasker, err error) {
	if taskerOption.err != nil {
		return nil, taskerOption.err
	}
	if at.IsZero() {
		return nil, errors.New("锚点时间不能为空")
	}
	if interval < 1 {
		interval = 1
	}

	tasker = &TaskerImpl{
		UUID:       _uuid.Must(_uuid.NewV7()),
		TaskType:   TaskTypeDaily,
		Name:       taskerOption.Name,
		Interval:   time.Duration(interval),
		At:         at,
		Timeout:    taskerOption.Timeout,
		Fn:         taskerOption.Fn,
		ErrHandler: taskerOption.ErrHandler,
		stop:       make(chan struct{}, 1),
	}
	return
}

func (*TaskerImpl) Hourly(interval int, at time.Time, taskerOption TaskerOption) (tasker Tasker, err error) {
	if taskerOption.err != nil {
		return nil, taskerOption.err
	}
	if at.IsZero() {
		return nil, errors.New("锚点时间不能为空")
	}
	if interval < 1 {
		interval = 1
	}

	tasker = &TaskerImpl{
		UUID:       _uuid.Must(_uuid.NewV7()),
		TaskType:   TaskTypeHourly,
		Name:       taskerOption.Name,
		Interval:   time.Duration(interval),
		At:         at,
		Timeout:    taskerOption.Timeout,
		Fn:         taskerOption.Fn,
		ErrHandler: taskerOption.ErrHandler,
		stop:       make(chan struct{}, 1),
	}
	return
}

func (*TaskerImpl) Minutely(interval int, at time.Time, taskerOption TaskerOption) (tasker Tasker, err error) {
	if taskerOption.err != nil {
		return nil, taskerOption.err
	}
	if at.IsZero() {
		return nil, errors.New("锚点时间不能为空")
	}
	if interval < 1 {
		interval = 1
	}

	tasker = &TaskerImpl{
		UUID:       _uuid.Must(_uuid.NewV7()),
		TaskType:   TaskTypeMinutely,
		Name:       taskerOption.Name,
		Interval:   time.Duration(interval),
		At:         at,
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
	case TaskTypeYearly, TaskTypeMonthly, TaskTypeWeekly, TaskTypeDaily, TaskTypeHourly, TaskTypeMinutely:
		return my.runSchedule(int(my.Interval))
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

// runSchedule 通用周期调度：计算下一个精确执行时刻，循环等待并执行
func (my *TaskerImpl) runSchedule(interval int) (err error) {
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
		next = my.advanceNext(next, loc, interval)
	}
}

// computeNext 根据任务类型和锚点，计算第一个未来执行时刻
func (my *TaskerImpl) computeNext(loc *time.Location) time.Time {
	var now = time.Now().In(loc)
	var h, m, s = my.At.Clock()

	switch my.TaskType {
	case TaskTypeYearly:
		var next = time.Date(now.Year(), my.At.Month(), my.At.Day(), h, m, s, 0, loc)
		if !next.After(now) {
			next = time.Date(now.Year()+1, my.At.Month(), my.At.Day(), h, m, s, 0, loc)
		}
		return next
	case TaskTypeMonthly:
		var next = time.Date(now.Year(), now.Month(), my.At.Day(), h, m, s, 0, loc)
		if !next.After(now) {
			next = time.Date(now.Year(), now.Month()+1, my.At.Day(), h, m, s, 0, loc)
		}
		return next
	case TaskTypeWeekly:
		var next = time.Date(now.Year(), now.Month(), now.Day(), h, m, s, 0, loc)
		var targetWeekday = my.At.Weekday()
		var daysAhead = (int(targetWeekday) - int(now.Weekday()) + 7) % 7
		if daysAhead == 0 && !next.After(now) {
			daysAhead = 7
		}
		next = next.AddDate(0, 0, daysAhead)
		return next
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
func (my *TaskerImpl) advanceNext(current time.Time, loc *time.Location, interval int) time.Time {
	var h, m, s = my.At.Clock()

	switch my.TaskType {
	case TaskTypeYearly:
		return time.Date(current.Year()+interval, my.At.Month(), my.At.Day(), h, m, s, 0, loc)
	case TaskTypeMonthly:
		return time.Date(current.Year(), current.Month()+time.Month(interval), my.At.Day(), h, m, s, 0, loc)
	case TaskTypeWeekly:
		return current.AddDate(0, 0, 7*interval)
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
