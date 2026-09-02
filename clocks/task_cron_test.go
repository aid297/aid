package clocks_test

import (
	_context "context"
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

	// 先注册到 clock，Stop 才能移除 cron entry；cron 调度器由 Boot 启动
	clock.AddTasker(tasker)
	clock.SetErrHandler(func(tasker _clocks.Tasker, err error) { t.Errorf("任务：%v 启动失败：%v", tasker, err) })
	clock.Boot(_context.Background())

	_time.Sleep(3200 * _time.Millisecond)

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

	// cron 未启动时 AddFunc 也能注册 entry，先同步调两次验证幂等
	if err := tasker.Begin(); err != nil {
		t.Fatalf("第一次 Begin 失败：%v", err)
	}
	if err := tasker.Begin(); err != nil {
		t.Fatalf("第二次 Begin 失败：%v", err)
	}

	// 启动调度器验证单 entry 只按秒触发一次
	clock.SetErrHandler(func(tasker _clocks.Tasker, err error) { t.Errorf("任务：%v 启动失败：%v", tasker, err) })
	clock.Boot(_context.Background())

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

// Test_CronGracefulShutdown 验证 Boot(ctx) 优雅退出：ctx 取消后 cron 调度器停止，任务不再触发
func Test_CronGracefulShutdown(t *_testing.T) {
	clock := _clocks.OnceClock()
	defer clock.Clean()

	var count int32
	tasker := _clocks.NewTaskCron().SetExpr("@every 1s")
	tasker.
		SetName("优雅退出测试").
		SetHandler(func(_ _clocks.Tasker) {
			_syncAtomic.AddInt32(&count, 1)
			t.Logf("触发第 %d 次: %s", _syncAtomic.LoadInt32(&count), _time.Now().Format("15:04:05.000"))
		})

	clock.AddTasker(tasker)

	ctx, cancel := _context.WithCancel(_context.Background())
	defer cancel()
	clock.Boot(ctx)

	_time.Sleep(2200 * _time.Millisecond)
	if got := _syncAtomic.LoadInt32(&count); got < 2 {
		t.Fatalf("取消前期望至少触发 2 次，实际 %d 次", got)
	}

	cancel()

	// 等待优雅停止生效（Boot 的 goroutine 会等待正在执行的回调结束后才退出）
	_time.Sleep(2200 * _time.Millisecond)
	paused := _syncAtomic.LoadInt32(&count)
	_time.Sleep(1500 * _time.Millisecond)
	if after := _syncAtomic.LoadInt32(&count); after != paused {
		t.Fatalf("ctx 取消后仍触发：%d -> %d", paused, after)
	}

	// 恢复全局 cron 调度器，避免影响后续用例
	clock.Boot(_context.Background())

	t.Logf("测试通过，取消前触发 %d 次，优雅退出后停止触发", paused)
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

	// cron 调度器由 Boot 启动
	clock.AddTasker(tasker)
	clock.SetErrHandler(func(tasker _clocks.Tasker, err error) { t.Errorf("任务：%v 启动失败：%v", tasker, err) })
	clock.Boot(_context.Background())

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

	t.Logf("测试通过，Close 前触发 %d 次，Close 后停止触发", paused)
}

// Test_BootGracefulShutdownNonCron 验证 Boot(ctx) 取消后非 cron 任务的行为：
// 正在执行的 handler 执行完才结束，未触发的任务立即注销不再生效
func Test_BootGracefulShutdownNonCron(t *_testing.T) {
	clock := _clocks.OnceClock()
	clock.SetErrHandler(func(tasker _clocks.Tasker, err error) { t.Errorf("任务：%v 执行失败：%v", tasker, err) })

	var fired, completed, onceFired int32

	// 每秒触发的周期任务，handler 执行 1.5 秒（跨越取消时刻，验证在跑任务收尾语义）
	cyc := _clocks.NewTaskCyclicitySecondly(1).
		SetName("周期任务").
		SetHandler(func(_ _clocks.Tasker) {
			_syncAtomic.AddInt32(&fired, 1)
			_time.Sleep(1500 * _time.Millisecond)
			_syncAtomic.AddInt32(&completed, 1)
		})

	// 远期一次性任务：取消后不应触发
	once := _clocks.NewTaskOnceAfter(10 * _time.Second).
		SetName("远期一次性任务").
		SetHandler(func(_ _clocks.Tasker) { _syncAtomic.AddInt32(&onceFired, 1) })

	clock.AddTasker(cyc, once)

	ctx, cancel := _context.WithCancel(_context.Background())
	clock.Boot(ctx)

	// 等待周期任务触发至少一次（handler 1.5s，在 2.2s 时一定还在跑）
	_time.Sleep(2200 * _time.Millisecond)
	if _syncAtomic.LoadInt32(&fired) == 0 {
		t.Fatal("取消前周期任务应已触发")
	}

	cancel()

	// 等待足够时间：在跑的 handler（最长 1.5s）应执行完，后续 tick 不再触发
	_time.Sleep(3 * _time.Second)
	f := _syncAtomic.LoadInt32(&fired)
	_time.Sleep(1500 * _time.Millisecond)

	// 正在执行的 handler 已执行完（fired == completed），后续不再触发，未触发的一次性任务未生效
	f2, c2 := _syncAtomic.LoadInt32(&fired), _syncAtomic.LoadInt32(&completed)
	if f2 != f {
		t.Fatalf("ctx 取消后周期任务仍触发：%d -> %d", f, f2)
	}
	if c2 != f2 {
		t.Fatalf("在跑 handler 未收尾：fired = %d, completed = %d", f2, c2)
	}
	if got := _syncAtomic.LoadInt32(&onceFired); got != 0 {
		t.Fatalf("远期一次性任务不应触发，实际触发 %d 次", got)
	}

	// 恢复全局 cron 调度器，避免影响后续用例（Clean 已由 Boot 的关停序列完成）
	clock.Boot(_context.Background())

	t.Logf("测试通过：取消后周期任务触发 %d 次全部收尾，远期任务未触发", f2)
}
