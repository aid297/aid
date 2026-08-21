package timers_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aid297/aid/v2/timers"
)

// errRecorder 收集 errHandler 收到的错误，并发安全
type errRecorder struct {
	mu   sync.Mutex
	errs []error
}

func (r *errRecorder) handler() timers.ErrHandler {
	return func(tasker timers.Tasker, err error) {
		r.mu.Lock()
		defer r.mu.Unlock()
		r.errs = append(r.errs, err)
	}
}

func (r *errRecorder) count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.errs)
}

func (r *errRecorder) last() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.errs) == 0 {
		return nil
	}
	return r.errs[len(r.errs)-1]
}

func newOnceTask(t *testing.T, fn timers.FN, interval time.Duration, errHandler timers.ErrHandler) timers.Tasker {
	t.Helper()
	tasker, err := timers.Task.New(timers.TaskTypeOnce, "test-once", interval, 0, fn, errHandler)
	if err != nil {
		t.Fatalf("New 失败：%v", err)
	}
	return tasker
}

// TestNew 验证 New 的参数校验
func TestNew(t *testing.T) {
	var noop = func() {}
	var handler = func(tasker timers.Tasker, err error) {}

	if _, err := timers.Task.New(timers.TaskTypeOnce, "t", 10*time.Millisecond, 0, noop, handler); err != nil {
		t.Fatalf("合法参数创建失败：%v", err)
	}

	if _, err := timers.Task.New(timers.TaskTypeTag(-1), "t", 10*time.Millisecond, 0, noop, handler); err == nil {
		t.Fatal("期望非法任务类型被拒绝")
	}
	if _, err := timers.Task.New(timers.TaskTypeOnce, "t", 0, 0, noop, handler); err == nil {
		t.Fatal("期望 interval <= 0 被拒绝")
	}
	if _, err := timers.Task.New(timers.TaskTypeOnce, "t", 10*time.Millisecond, 0, nil, handler); err == nil {
		t.Fatal("期望 nil fn 被拒绝")
	}
	if _, err := timers.Task.New(timers.TaskTypeOnce, "t", 10*time.Millisecond, 0, noop, nil); err == nil {
		t.Fatal("期望 nil errHandler 被拒绝")
	}
}

// TestOnceTaskSuccess 验证一次性任务正常执行并返回
func TestOnceTaskSuccess(t *testing.T) {
	var count int64
	var rec = &errRecorder{}
	var start = time.Now()
	tasker := newOnceTask(t, func() { atomic.AddInt64(&count, 1) }, 20*time.Millisecond, rec.handler())

	if err := tasker.Start(); err != nil {
		t.Fatalf("Start 失败：%v", err)
	}

	if got := atomic.LoadInt64(&count); got != 1 {
		t.Fatalf("期望执行 1 次，实际 %d 次", got)
	}
	if time.Since(start) < 20*time.Millisecond {
		t.Fatal("任务未等到 interval 就执行了")
	}
	if rec.count() != 0 {
		t.Fatalf("期望无错误上报，实际 %d 个", rec.count())
	}
}

// TestOnceTaskPanic 验证 fn panic 被捕获并上报 errHandler
func TestOnceTaskPanic(t *testing.T) {
	var rec = &errRecorder{}
	tasker := newOnceTask(t, func() { panic("模拟任务异常") }, 10*time.Millisecond, rec.handler())

	if err := tasker.Start(); err == nil {
		t.Fatal("期望 panic 时 Start 返回错误")
	}
	if rec.count() != 1 {
		t.Fatalf("期望 errHandler 收到 1 个错误，实际 %d 个", rec.count())
	}
}

// TestOnceTaskStop 验证执行前 Stop，fn 不执行且不上报错误
func TestOnceTaskStop(t *testing.T) {
	var count int64
	var rec = &errRecorder{}
	tasker := newOnceTask(t, func() { atomic.AddInt64(&count, 1) }, 50*time.Millisecond, rec.handler())

	result := make(chan error, 1)
	go func() { result <- tasker.Start() }()

	time.Sleep(10 * time.Millisecond)
	if err := tasker.Stop(); err != nil {
		t.Fatalf("Stop 失败：%v", err)
	}

	select {
	case err := <-result:
		if err != nil {
			t.Fatalf("期望 Stop 后 Start 返回 nil（错误走 errHandler），实际：%v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Stop 后 Start 未及时返回")
	}

	time.Sleep(80 * time.Millisecond)
	if got := atomic.LoadInt64(&count); got != 0 {
		t.Fatalf("停止后 fn 不应执行，实际执行 %d 次", got)
	}
	if rec.count() != 0 {
		t.Fatalf("手动停止不算错误，期望无上报，实际 %d 个", rec.count())
	}
}

// TestOnceTaskTimeout 验证超时时经 errHandler 上报且 fn 不执行
func TestOnceTaskTimeout(t *testing.T) {
	var count int64
	var rec = &errRecorder{}
	tasker, err := timers.Task.New(
		timers.TaskTypeOnce, "test-once-timeout", 50*time.Millisecond, 20*time.Millisecond,
		func() { atomic.AddInt64(&count, 1) },
		rec.handler(),
	)
	if err != nil {
		t.Fatalf("New 失败：%v", err)
	}

	if err = tasker.Start(); err != nil {
		t.Fatalf("期望超时后 Start 返回 nil（错误走 errHandler），实际：%v", err)
	}
	if rec.count() != 1 {
		t.Fatalf("期望 errHandler 收到超时错误，实际 %d 个", rec.count())
	}

	time.Sleep(80 * time.Millisecond)
	if got := atomic.LoadInt64(&count); got != 0 {
		t.Fatalf("超时退出后 fn 不应执行，实际执行 %d 次", got)
	}
}

// TestUnsupportedTaskType 验证暂不支持的任务类型返回错误
func TestUnsupportedTaskType(t *testing.T) {
	var rec = &errRecorder{}
	tasker, err := timers.Task.New(
		timers.TaskTypeWeekly, "test-weekly", 10*time.Millisecond, 0,
		func() {},
		rec.handler(),
	)
	if err != nil {
		t.Fatalf("New 失败：%v", err)
	}

	if err = tasker.Start(); err == nil {
		t.Fatal("期望不支持的任务类型 Start 失败")
	}
}

// TestStopIdempotent 验证重复调用 Stop 不会 panic
func TestStopIdempotent(t *testing.T) {
	var rec = &errRecorder{}
	tasker, err := timers.Task.New(
		timers.TaskTypeOnce, "test-stop-idempotent", 10*time.Millisecond, 0,
		func() {},
		rec.handler(),
	)
	if err != nil {
		t.Fatalf("New 失败：%v", err)
	}

	for i := 0; i < 3; i++ {
		if err = tasker.Stop(); err != nil {
			t.Fatalf("第 %d 次 Stop 失败：%v", i+1, err)
		}
	}
}

// TestDailyTaskLoop 验证按天循环任务按周期重复执行，Stop 后退出且不再执行
func TestDailyTaskLoop(t *testing.T) {
	var count int64
	var rec = &errRecorder{}
	tasker, err := timers.Task.New(
		timers.TaskTypeDaily, "test-daily", 20*time.Millisecond, 0,
		func() { atomic.AddInt64(&count, 1) },
		rec.handler(),
	)
	if err != nil {
		t.Fatalf("New 失败：%v", err)
	}

	result := make(chan error, 1)
	go func() { result <- tasker.Start() }()

	time.Sleep(70 * time.Millisecond)
	if got := atomic.LoadInt64(&count); got < 2 {
		t.Fatalf("期望至少执行 2 次，实际 %d 次", got)
	}

	if err = tasker.Stop(); err != nil {
		t.Fatalf("Stop 失败：%v", err)
	}
	select {
	case err = <-result:
		if err != nil {
			t.Fatalf("期望 Stop 后 Start 返回 nil，实际：%v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Stop 后 Start 未及时返回")
	}

	stopped := atomic.LoadInt64(&count)
	time.Sleep(60 * time.Millisecond)
	if got := atomic.LoadInt64(&count); got != stopped {
		t.Fatalf("停止后不应再执行：停止时 %d 次，当前 %d 次", stopped, got)
	}
}

// TestDailyTaskPanicContinues 验证单次 panic 被恢复、上报且不中断循环调度
func TestDailyTaskPanicContinues(t *testing.T) {
	var count int64
	var rec = &errRecorder{}
	tasker, err := timers.Task.New(
		timers.TaskTypeDaily, "test-daily-panic", 20*time.Millisecond, 0,
		func() { atomic.AddInt64(&count, 1); panic("模拟任务异常") },
		rec.handler(),
	)
	if err != nil {
		t.Fatalf("New 失败：%v", err)
	}

	go tasker.Start()
	time.Sleep(70 * time.Millisecond)
	tasker.Stop()

	if got := atomic.LoadInt64(&count); got < 2 {
		t.Fatalf("期望 panic 后继续调度，至少执行 2 次，实际 %d 次", got)
	}
	if rec.count() < 2 {
		t.Fatalf("期望每次 panic 都上报 errHandler，实际 %d 个", rec.count())
	}
}

// TestDailyTaskTimeout 验证循环任务等待超时时经 errHandler 上报且任务不再执行
func TestDailyTaskTimeout(t *testing.T) {
	var count int64
	var rec = &errRecorder{}
	tasker, err := timers.Task.New(
		timers.TaskTypeDaily, "test-daily-timeout", 50*time.Millisecond, 20*time.Millisecond,
		func() { atomic.AddInt64(&count, 1) },
		rec.handler(),
	)
	if err != nil {
		t.Fatalf("New 失败：%v", err)
	}

	if err = tasker.Start(); err != nil {
		t.Fatalf("期望超时后 Start 返回 nil，实际：%v", err)
	}
	if rec.count() != 1 {
		t.Fatalf("期望 errHandler 收到超时错误，实际 %d 个", rec.count())
	}

	time.Sleep(80 * time.Millisecond)
	if got := atomic.LoadInt64(&count); got != 0 {
		t.Fatalf("超时退出后任务不应执行，实际执行 %d 次", got)
	}
}
