package clocks_test

import (
	_sync "sync"
	_syncAtomic "sync/atomic"
	_testing "testing"
	_time "time"

	_clocks "github.com/aid297/aid/v2/clocks"
)

// Test_CronAttrs 验证属性链式设置与读取
func Test_CronAttrs(t *_testing.T) {
	tasker1 := _clocks.NewTaskCron()
	tasker2 := _clocks.NewTaskCron()

	if tasker1.UUID() == tasker2.UUID() {
		t.Fatal("两个任务的 UUID 不应相同")
	}

	var handler _clocks.TaskHandler = func(_ _clocks.Tasker) {}
	tasker1.
		SetExpr("@every 1s").
		SetName("cron任务").
		SetTimeout(30 * _time.Second).
		SetHandler(handler)

	if tasker1.Name() != "cron任务" {
		t.Fatalf("Name() = %q, want %q", tasker1.Name(), "cron任务")
	}
	if tasker1.Timeout() != 30*_time.Second {
		t.Fatalf("Timeout() = %s, want %s", tasker1.Timeout(), 30*_time.Second)
	}
	if tasker1.Handler() == nil {
		t.Fatal("Handler() 不应为空")
	}
	if tasker1.Expr() != "@every 1s" {
		t.Fatalf("Expr() = %q, want %q", tasker1.Expr(), "@every 1s")
	}
	if tasker1.String() == "" {
		t.Fatal("String() 不应为空")
	}

	// cron 任务不涉及 immediately 语义，三个方法均为空实现，调用不应 panic
	tasker1.SetImmediately(false).EnableImmediately().DisableImmediately()
	if !tasker1.Immediately() {
		t.Fatal("Immediately() 恒为 true")
	}

	// EntryID 初始应为 0
	if tasker1.EntryID() != 0 {
		t.Fatalf("EntryID() 初始应为 0，实际 %d", tasker1.EntryID())
	}

	t.Log("测试通过，属性设置与读取均正常")
}

// Test_CronNoHandler 验证未设置回调时 Begin 返回错误
func Test_CronNoHandler(t *_testing.T) {
	tasker := _clocks.NewTaskCron().SetExpr("@every 1s")

	err := tasker.Begin()
	if err == nil {
		t.Fatal("期望返回错误，但未返回")
	}
	t.Logf("符合预期：%v", err)
}

// Test_CronNoExpr 验证未设置表达式时 Begin 返回错误
func Test_CronNoExpr(t *_testing.T) {
	tasker := _clocks.NewTaskCron().SetHandler(func(_ _clocks.Tasker) {})

	err := tasker.Begin()
	if err == nil {
		t.Fatal("期望返回错误，但未返回")
	}
	t.Logf("符合预期：%v", err)
}

// Test_CronInvalidExpr 验证非法 cron 表达式 Begin 返回解析错误
func Test_CronInvalidExpr(t *_testing.T) {
	tasker := _clocks.NewTaskCron().
		SetExpr("invalid-expr").
		SetHandler(func(_ _clocks.Tasker) {})

	err := tasker.Begin()
	if err == nil {
		t.Fatal("期望返回错误，但未返回")
	}
	t.Logf("符合预期：%v", err)
}

// Test_CronEverySecond 验证 @every 1s 表达式按秒触发，且 Stop 后不再触发
func Test_CronEverySecond(t *_testing.T) {
	clock := _clocks.OnceClock()
	defer clock.Clean()

	var count int32
	tasker := _clocks.NewTaskCron().SetExpr("@every 1s")
	tasker.
		SetName("每秒执行").
		SetHandler(func(_ _clocks.Tasker) {
			_syncAtomic.AddInt32(&count, 1)
			t.Logf("触发第 %d 次: %s", _syncAtomic.LoadInt32(&count), _time.Now().Format("15:04:05.000"))
		})

	// 先注册到 clock，Stop 才能移除 cron entry
	clock.AddTasker(tasker)

	errCh := make(chan error, 1)
	go func() { errCh <- tasker.Begin() }()

	_time.Sleep(3200 * _time.Millisecond)

	if err := <-errCh; err != nil {
		t.Fatalf("任务启动失败：%v", err)
	}

	if tasker.EntryID() == 0 {
		t.Fatal("Begin 后 EntryID 不应为 0")
	}

	got := _syncAtomic.LoadInt32(&count)
	if got < 2 {
		t.Fatalf("期望 3 秒内至少触发 2 次，实际 %d 次", got)
	}

	// Stop 之后不应再触发
	tasker.Stop()
	paused := _syncAtomic.LoadInt32(&count)
	_time.Sleep(2 * _time.Second)
	if after := _syncAtomic.LoadInt32(&count); after != paused {
		t.Fatalf("Stop 之后仍触发：%d -> %d", paused, after)
	}

	t.Logf("测试通过，共触发 %d 次，Stop 后停止触发", got)
}

// Test_CronBeginIdempotent 验证重复 Begin 不会重复注册 cron entry
func Test_CronBeginIdempotent(t *_testing.T) {
	clock := _clocks.OnceClock()
	defer clock.Clean()

	var count int32
	tasker := _clocks.NewTaskCron().
		SetExpr("@every 1s").
		SetName("幂等测试").
		SetHandler(func(_ _clocks.Tasker) { _syncAtomic.AddInt32(&count, 1) })

	clock.AddTasker(tasker)

	if err := tasker.Begin(); err != nil {
		t.Fatalf("第一次 Begin 失败：%v", err)
	}
	if err := tasker.Begin(); err != nil {
		t.Fatalf("第二次 Begin 失败：%v", err)
	}

	_time.Sleep(3200 * _time.Millisecond)

	// 若重复注册，3 秒内会触发约 6 次；正常只应触发 2-4 次
	got := _syncAtomic.LoadInt32(&count)
	if got < 2 || got > 4 {
		t.Fatalf("重复 Begin 后触发次数异常：%d 次（期望 2-4 次）", got)
	}

	tasker.Stop()
	t.Logf("测试通过，触发 %d 次，未重复注册", got)
}

// Test_CronErrHandlerViaClock 验证非法表达式经 AddTaskerAndBegin 启动时，错误被 errHandler 捕获
func Test_CronErrHandlerViaClock(t *_testing.T) {
	clock := _clocks.OnceClock()

	var wg _sync.WaitGroup
	wg.Add(1)
	errCh := make(chan error, 1)
	clock.SetErrHandler(func(_ _clocks.Tasker, err error) {
		defer wg.Done()
		errCh <- err
	})
	// 测试结束后恢复默认处理，避免影响其他用例
	defer clock.SetErrHandler(func(_ _clocks.Tasker, _ error) {})

	tasker := _clocks.NewTaskCron().
		SetExpr("bad-expr").
		SetName("非法表达式").
		SetHandler(func(_ _clocks.Tasker) {})

	clock.AddTaskerAndBegin(tasker)

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatal("errHandler 收到 nil 错误")
		}
		t.Logf("错误被正确上报：%v", err)
	case <-_time.After(3 * _time.Second):
		t.Fatal("errHandler 未收到错误")
	}

	clock.DeleteTasker(tasker.UUID())
	wg.Wait()
}

// Test_CronCloseByClock 验证通过 clock.Close 停止任务后不再触发，且任务被移除
func Test_CronCloseByClock(t *_testing.T) {
	clock := _clocks.OnceClock()
	defer clock.Clean()

	var count int32
	tasker := _clocks.NewTaskCron().
		SetExpr("@every 1s").
		SetName("Close测试").
		SetHandler(func(_ _clocks.Tasker) { _syncAtomic.AddInt32(&count, 1) })

	errCh := make(chan error, 1)
	go func() { errCh <- tasker.Begin() }()

	clock.AddTasker(tasker)

	if got := clock.Tasker(tasker.UUID()); got == nil {
		t.Fatal("Tasker() 查询失败，未找到已注册任务")
	}

	_time.Sleep(2200 * _time.Millisecond)

	clock.Close(tasker.UUID())

	paused := _syncAtomic.LoadInt32(&count)
	_time.Sleep(2 * _time.Second)
	if after := _syncAtomic.LoadInt32(&count); after != paused {
		t.Fatalf("Close 之后仍触发：%d -> %d", paused, after)
	}

	if clock.Tasker(tasker.UUID()) != nil {
		t.Fatal("Close 后任务应被移除")
	}

	if err := <-errCh; err != nil {
		t.Fatalf("任务启动失败：%v", err)
	}

	t.Logf("测试通过，Close 前触发 %d 次，Close 后停止触发", paused)
}
