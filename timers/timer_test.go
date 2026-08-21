package timers_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aid297/aid/v2/timers"
	_uuid "github.com/google/uuid"
)

// timerErrRecorder 收集 errHandler 收到的错误，并发安全
type timerErrRecorder struct {
	mu   sync.Mutex
	errs []error
}

func (r *timerErrRecorder) handler() timers.ErrHandler {
	return func(tasker timers.Tasker, err error) {
		r.mu.Lock()
		defer r.mu.Unlock()
		r.errs = append(r.errs, err)
	}
}

func (r *timerErrRecorder) count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.errs)
}

func newTimerOnceTask(t *testing.T, fn timers.FN, errHandler timers.ErrHandler) timers.Tasker {
	t.Helper()
	tasker, err := timers.Task.New(timers.TaskTypeOnce, "timer-test-once", 10*time.Millisecond, 0, fn, errHandler)
	if err != nil {
		t.Fatalf("New 失败：%v", err)
	}
	return tasker
}

func newTimerDailyTask(t *testing.T, name string, fn timers.FN, errHandler timers.ErrHandler) timers.Tasker {
	t.Helper()
	tasker, err := timers.Task.New(timers.TaskTypeDaily, name, 20*time.Millisecond, 0, fn, errHandler)
	if err != nil {
		t.Fatalf("New 失败：%v", err)
	}
	return tasker
}

// TestTimerOnce 验证单例非空且多次获取为同一实例
func TestTimerOnce(t *testing.T) {
	var ins1 = timers.Time.Once()
	var ins2 = timers.Time.Once()

	if ins1 == nil || ins2 == nil {
		t.Fatal("单例不应为 nil")
	}
	if ins1 != ins2 {
		t.Fatal("多次 Once 应返回同一实例")
	}
}

// TestTimerAddAndGetTask 验证任务登记与查询，nil 任务被安全忽略
func TestTimerAddAndGetTask(t *testing.T) {
	var ins = timers.Time.Once()

	var tasker = newTimerOnceTask(t, func() {}, func(tasker timers.Tasker, err error) {})
	ins.AddTask(tasker)

	if got := ins.GetTask(tasker.GetUUID()); got == nil {
		t.Fatal("已登记任务应能查询到")
	}
	if got := ins.GetTask(_uuid.Must(_uuid.NewV7())); got != nil {
		t.Fatal("未登记任务应返回 nil")
	}

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("MustGetTask 查询不存在的任务应 panic")
			}
		}()
		ins.MustGetTask(_uuid.Must(_uuid.NewV7()))
	}()

	// nil 任务不应 panic
	ins.AddTask(nil)
	ins.AddTaskAndStart(nil)
}

// TestTimerStartRunsAllTasks 验证 Start 启动全部已登记任务，重复 Start 不再执行
func TestTimerStartRunsAllTasks(t *testing.T) {
	var ins = timers.Time.Once()

	var count1, count2 int64
	ins.AddTask(newTimerOnceTask(t, func() { atomic.AddInt64(&count1, 1) }, func(tasker timers.Tasker, err error) {}))
	ins.AddTask(newTimerOnceTask(t, func() { atomic.AddInt64(&count2, 1) }, func(tasker timers.Tasker, err error) {}))

	ins.Start()
	time.Sleep(100 * time.Millisecond)

	if got := atomic.LoadInt64(&count1); got != 1 {
		t.Fatalf("任务 1 期望执行 1 次，实际 %d 次", got)
	}
	if got := atomic.LoadInt64(&count2); got != 1 {
		t.Fatalf("任务 2 期望执行 1 次，实际 %d 次", got)
	}

	// 重复 Start 被 running 挡住，一次性任务不会再次执行
	ins.Start()
	time.Sleep(100 * time.Millisecond)
	if got := atomic.LoadInt64(&count1); got != 1 {
		t.Fatalf("重复 Start 后任务 1 不应再次执行，实际 %d 次", got)
	}
}

// TestTimerAddTaskWhileRunning 验证运行中 AddTask 不会执行，AddTaskAndStart 立即执行
func TestTimerAddTaskWhileRunning(t *testing.T) {
	var ins = timers.Time.Once()

	var lateCount, dynCount int64
	ins.AddTask(newTimerOnceTask(t, func() { atomic.AddInt64(&lateCount, 1) }, func(tasker timers.Tasker, err error) {}))
	ins.AddTaskAndStart(newTimerOnceTask(t, func() { atomic.AddInt64(&dynCount, 1) }, func(tasker timers.Tasker, err error) {}))

	time.Sleep(100 * time.Millisecond)

	if got := atomic.LoadInt64(&lateCount); got != 0 {
		t.Fatalf("运行中 AddTask 的任务不应自动执行，实际 %d 次", got)
	}
	if got := atomic.LoadInt64(&dynCount); got != 1 {
		t.Fatalf("AddTaskAndStart 的任务应立即执行，实际 %d 次", got)
	}
}

// TestTimerStop 验证 Stop 停止全部任务，手动停止不上报错误，回调内调锁方法不死锁
func TestTimerStop(t *testing.T) {
	var ins = timers.Time.Once()

	var count int64
	var rec = &timerErrRecorder{}
	var tasker = newTimerDailyTask(t, "timer-test-stop", func() { atomic.AddInt64(&count, 1) },
		func(tasker timers.Tasker, err error) {
			// 回调内再调锁方法，验证不死锁
			ins.GetTask(tasker.GetUUID())
			rec.handler()(tasker, err)
		})
	ins.AddTaskAndStart(tasker)

	time.Sleep(70 * time.Millisecond)
	if got := atomic.LoadInt64(&count); got < 2 {
		t.Fatalf("任务应先正常运行，实际执行 %d 次", got)
	}

	ins.Stop()
	time.Sleep(100 * time.Millisecond)

	stopped := atomic.LoadInt64(&count)
	time.Sleep(60 * time.Millisecond)
	if got := atomic.LoadInt64(&count); got != stopped {
		t.Fatalf("Stop 后任务不应再执行：停止时 %d 次，当前 %d 次", stopped, got)
	}
	if rec.count() != 0 {
		t.Fatalf("手动停止不算错误，期望无上报，实际 %d 个", rec.count())
	}
}

// TestTimerRestartAfterStop 验证 Stop 复位 running 后可再次 Start
func TestTimerRestartAfterStop(t *testing.T) {
	var ins = timers.Time.Once()

	var count int64
	ins.AddTask(newTimerOnceTask(t, func() { atomic.AddInt64(&count, 1) }, func(tasker timers.Tasker, err error) {}))

	ins.Start()
	time.Sleep(100 * time.Millisecond)
	if got := atomic.LoadInt64(&count); got != 1 {
		t.Fatalf("重启后任务期望执行 1 次，实际 %d 次", got)
	}
}

// TestTimerTaskErrHandlerPanicProtection 验证任务 errHandler panic 不拖垮进程，timer 仍可继续使用
func TestTimerTaskErrHandlerPanicProtection(t *testing.T) {
	var ins = timers.Time.Once()

	// fn panic 触发 errHandler，errHandler 再 panic，应由 runTask 的 recover 兜底
	var badHandler = func(tasker timers.Tasker, err error) { panic("handler 异常") }
	ins.AddTaskAndStart(newTimerOnceTask(t, func() { panic("模拟失败") }, badHandler))

	time.Sleep(100 * time.Millisecond)

	// 进程存活且 timer 可继续使用
	var count int64
	ins.AddTaskAndStart(newTimerOnceTask(t, func() { atomic.AddInt64(&count, 1) }, func(tasker timers.Tasker, err error) {}))

	time.Sleep(100 * time.Millisecond)
	if got := atomic.LoadInt64(&count); got != 1 {
		t.Fatalf("handler panic 后 timer 应仍可正常运行任务，实际执行 %d 次", got)
	}
}
