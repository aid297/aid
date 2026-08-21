package timers

import (
	"errors"
	"fmt"
	"sync"
	"time"

	_uuid "github.com/google/uuid"
)

const (
	TaskTypeOnce TaskTypeTag = iota
	TaskTypeCyclicity
)

var (
	TaskIntervalDaily  = time.Duration(24) * time.Hour
	TaskIntervalWeekly = time.Duration(7) * TaskIntervalDaily
)

type (
	TaskTypeTag int

	Tasker interface {
		New(
			taskType TaskTypeTag,
			name string,
			interval time.Duration,
			timeout time.Duration,
			fn FN,
			errHandler ErrHandler,
		) (tasker Tasker, err error)
		GetUUID() _uuid.UUID
		Start() (err error)
		Stop() (err error)
	}

	TaskerImpl struct {
		uuid       _uuid.UUID
		name       string
		taskType   TaskTypeTag
		fn         FN
		interval   time.Duration
		stop       chan struct{}
		stopOnce   sync.Once
		timeout    time.Duration
		errHandler ErrHandler
	}

	FN func()

	ErrHandler func(tasker Tasker, err error)
)

func (*TaskerImpl) New(
	taskType TaskTypeTag,
	name string,
	interval time.Duration,
	timeout time.Duration,
	fn FN,
	errHandler ErrHandler,
) (tasker Tasker, err error) {
	if taskType < TaskTypeOnce || taskType > TaskTypeCyclicity {
		return nil, errors.New("定时任务类型错误")
	}
	if interval <= 0 {
		return nil, errors.New("时间错误")
	}
	if fn == nil {
		return nil, errors.New("回调方法不能为空")
	}
	if errHandler == nil {
		return nil, errors.New("错误处理方法不能为空")
	}

	tasker = &TaskerImpl{
		uuid:       _uuid.Must(_uuid.NewV7()),
		taskType:   taskType,
		name:       name,
		interval:   interval,
		timeout:    timeout,
		fn:         fn,
		errHandler: errHandler,
		stop:       make(chan struct{}, 1),
	}
	return
}

func (my *TaskerImpl) GetUUID() _uuid.UUID { return my.uuid }

func (my *TaskerImpl) Start() (err error) {
	if my.fn == nil {
		return fmt.Errorf("定时任务执行失败，没有需要执行的方法")
	}

	switch my.taskType {
	case TaskTypeOnce:
		var c = make(chan error, 1)
		var timer = time.AfterFunc(my.interval, func() {
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
			my.fn()
			c <- nil // 缓冲 channel，不会阻塞泄漏
		})

		var timeoutCh <-chan time.Time
		if my.timeout > 0 {
			timeoutCh = time.After(my.timeout) // timeout <= 0 时 timeoutCh 为 nil，该 case 永不触发
		}
		select {
		case err = <-c: // 执行完成，失败时上报 errHandler
			if err != nil {
				my.errHandler(my, err)
			}
		case <-my.stop: // 手动停止，不算错误
			timer.Stop()
			return
		case <-timeoutCh: // 等待超时
			timer.Stop()
			my.errHandler(my, fmt.Errorf("定时任务执行超时：%s", my.timeout))
			return
		}
		return
	case TaskTypeCyclicity:
		var ticker = time.NewTicker(my.interval)
		defer ticker.Stop()

		var timeoutCh <-chan time.Time
		if my.timeout > 0 {
			timeoutCh = time.After(my.timeout) // timeout <= 0 时 timeoutCh 为 nil，该 case 永不触发
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
					my.errHandler(my, fnErr)
				}
			case <-my.stop: // 手动停止，不算错误
				return
			case <-timeoutCh: // 等待超时
				my.errHandler(my, fmt.Errorf("定时任务执行超时：%s", my.timeout))
				return
			}
		}
	default:
		return fmt.Errorf("不支持的定时任务类型 %d", my.taskType)
	}
}

// safeFn 带 panic 保护地执行回调，供循环任务使用，单次 panic 不中断调度
func (my *TaskerImpl) safeFn() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("定时任务执行失败：%v", r)
		}
	}()

	my.fn()
	return
}

func (my *TaskerImpl) Stop() (err error) { my.stopOnce.Do(func() { close(my.stop) }); return }
