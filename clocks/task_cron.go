package clocks

import (
	_errors "errors"
	_fmt "fmt"
	_sync "sync"
	_time "time"

	_uuid "github.com/google/uuid"
	_cron "github.com/robfig/cron/v3"
)

type (
	TaskCronImpl struct {
		err     error
		arg     any
		handler TaskHandler
		name    string
		expr    string
		mu      _sync.Mutex
		entryID _cron.EntryID
		uuid    _uuid.UUID
		timeout _time.Duration
	}
)

func NewTaskCron() *TaskCronImpl {
	return &TaskCronImpl{uuid: _uuid.Must(_uuid.NewV7()), timeout: defaultTimeout}
}

func (my *TaskCronImpl) String() string {
	return _fmt.Sprintf("uuid: %s, name: %s, timeout: %s", my.uuid, my.name, my.timeout)
}

func (my *TaskCronImpl) UUID() _uuid.UUID { return my.uuid }

func (my *TaskCronImpl) SetName(name string) Tasker { my.name = name; return my }

func (my *TaskCronImpl) Name() string { return my.name }

func (my *TaskCronImpl) SetTimeout(timeout _time.Duration) Tasker { my.timeout = timeout; return my }

func (my *TaskCronImpl) Timeout() _time.Duration { return my.timeout }

func (my *TaskCronImpl) SetHandler(fn TaskHandler) Tasker { my.handler = fn; return my }

func (my *TaskCronImpl) Handler() TaskHandler { return my.handler }

func (my *TaskCronImpl) SetImmediately(_ bool) Tasker { return my }

func (my *TaskCronImpl) EnableImmediately() Tasker { return my }

func (my *TaskCronImpl) DisableImmediately() Tasker { return my }

func (my *TaskCronImpl) Immediately() bool { return true }

func (my *TaskCronImpl) SetExpr(expr string) *TaskCronImpl { my.expr = expr; return my }

func (my *TaskCronImpl) Expr() string { return my.expr }

func (my *TaskCronImpl) SetEntryID(entryID _cron.EntryID) Tasker { my.entryID = entryID; return my }

func (my *TaskCronImpl) EntryID() _cron.EntryID { return my.entryID }

func (my *TaskCronImpl) Begin() (err error) {
	if my.handler == nil {
		return _errors.New("回调方法为空")
	}

	if my.expr == "" {
		return _errors.New("cron 表达式为空")
	}

	my.mu.Lock()
	defer my.mu.Unlock()

	if my.entryID == 0 {
		my.entryID, err = clockIns.cron.AddFunc(my.expr, func() { my.handler(my) })
	}

	return
}

func (my *TaskCronImpl) Stop() Tasker {
	if _, ok := clockIns.taskers[my.uuid]; ok {
		my.mu.Lock()
		entryID := my.entryID
		my.mu.Unlock()

		if entryID != 0 {
			clockIns.cron.Remove(entryID)
		}

		delete(clockIns.taskers, my.uuid)
	}

	return my
}
