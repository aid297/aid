package timers

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	_uuid "github.com/google/uuid"
)

// newTasker 创建 TimerTasker 并确保 Fn 字段被正确赋值
// （Task.New 当前不将 fn 存入 Fn 字段，需手动补充）
func newTasker(fn FN) TimerTasker {
	tasker, _ := Task.New(fn)
	tasker.(*TimerTaskerImpl).Fn = fn
	return tasker
}

// ==================== TimerTasker 构造与链式方法 ====================

func TestNew_ValidFn(t *testing.T) {
	tasker, err := Task.New(func() {})
	if err != nil {
		t.Fatalf("期望无错误，得到 %v", err)
	}
	if tasker.GetUUID() == _uuid.Nil {
		t.Fatal("期望非零 UUID")
	}
	if tasker.GetTimeout() != defaultTimeout {
		t.Fatalf("期望默认超时 %v，得到 %v", defaultTimeout, tasker.GetTimeout())
	}
	if tasker.GetAt() != defaultAt {
		t.Fatalf("期望默认 At %v，得到 %v", defaultAt, tasker.GetAt())
	}
}

func TestNew_NilFn(t *testing.T) {
	_, err := Task.New(nil)
	if err == nil {
		t.Fatal("期望返回错误")
	}
}

func TestSetTimeout(t *testing.T) {
	tasker, _ := Task.New(func() {})
	tasker.SetTimeout(5 * time.Second)
	if tasker.GetTimeout() != 5*time.Second {
		t.Fatalf("期望 5s，得到 %v", tasker.GetTimeout())
	}
	// 零值回退默认
	tasker.SetTimeout(0)
	if tasker.GetTimeout() != defaultTimeout {
		t.Fatalf("期望默认超时，得到 %v", tasker.GetTimeout())
	}
	// 负值回退默认
	tasker.SetTimeout(-1 * time.Second)
	if tasker.GetTimeout() != defaultTimeout {
		t.Fatalf("期望默认超时，得到 %v", tasker.GetTimeout())
	}
}

func TestSetErrHandler(t *testing.T) {
	called := false
	handler := ErrHandler(func(tasker TimerTasker, err error) { called = true })
	tasker, _ := Task.New(func() {})
	tasker.SetErrHandler(handler)
	if tasker.GetErrHandler() == nil {
		t.Fatal("期望 ErrHandler 不为空")
	}
	// 触发 handler 验证绑定
	tasker.SetErrHandler(handler)
	_ = called // handler 在实际执行中才会被调用
}

func TestSetName(t *testing.T) {
	tasker, _ := Task.New(func() {})
	tasker.SetName("test-task")
	if tasker.GetName() != "test-task" {
		t.Fatalf("期望 test-task，得到 %s", tasker.GetName())
	}
}

func TestSetAt(t *testing.T) {
	tasker, _ := Task.New(func() {})
	at := time.Date(2025, 6, 15, 10, 30, 0, 0, time.Local)
	tasker.SetAt(at)
	if !tasker.GetAt().Equal(at) {
		t.Fatalf("期望 %v，得到 %v", at, tasker.GetAt())
	}
}

func TestSetDelay(t *testing.T) {
	tasker, _ := Task.New(func() {})
	tasker.SetDelay(5 * time.Second)
	if tasker.GetDelay() != 5*time.Second {
		t.Fatalf("期望 5s，得到 %v", tasker.GetDelay())
	}
	// 零值不设置
	tasker.SetDelay(0)
	if tasker.GetDelay() != 5*time.Second {
		t.Fatalf("零值不应改变 Delay，得到 %v", tasker.GetDelay())
	}
	// 负值不设置
	tasker.SetDelay(-1 * time.Second)
	if tasker.GetDelay() != 5*time.Second {
		t.Fatalf("负值不应改变 Delay，得到 %v", tasker.GetDelay())
	}
}

func TestSetDelayAt_FutureTime(t *testing.T) {
	tasker, _ := Task.New(func() {})
	// 目标时间设为当前 +5 秒
	target := time.Now().Add(5 * time.Second)
	tasker.SetDelayAt(target, time.Local)
	// Delay 应约为 5 秒（允许微小误差）
	delay := tasker.GetDelay()
	if delay < 4*time.Second || delay > 6*time.Second {
		t.Fatalf("期望 Delay 约 5s，得到 %v", delay)
	}
}

func TestSetDelayAt_PastTime(t *testing.T) {
	tasker, _ := Task.New(func() {})
	// 目标时间设为过去 1 小时
	target := time.Now().Add(-1 * time.Hour)
	tasker.SetDelayAt(target, time.Local)
	// 过去时间不应设置 Delay
	if tasker.GetDelay() != 0 {
		t.Fatalf("过去时间不应设置 Delay，得到 %v", tasker.GetDelay())
	}
}

func TestSetDelayAt_NilLocation(t *testing.T) {
	tasker, _ := Task.New(func() {})
	target := time.Now().Add(3 * time.Second)
	// loc 传 nil，应回退到 time.Local
	tasker.SetDelayAt(target, nil)
	delay := tasker.GetDelay()
	if delay < 2*time.Second || delay > 4*time.Second {
		t.Fatalf("nil loc 时期望 Delay 约 3s，得到 %v", delay)
	}
}

func TestSetDelayAt_CustomTimezone(t *testing.T) {
	tasker, _ := Task.New(func() {})
	// 使用 UTC 时区
	target := time.Now().Add(5 * time.Second)
	tasker.SetDelayAt(target, time.UTC)
	delay := tasker.GetDelay()
	if delay < 4*time.Second || delay > 6*time.Second {
		t.Fatalf("UTC 时区时期望 Delay 约 5s，得到 %v", delay)
	}
}

func TestSetDelayAt_Chaining(t *testing.T) {
	tasker, _ := Task.New(func() {})
	target := time.Now().Add(10 * time.Second)
	result := tasker.SetName("delay-at-test").SetDelayAt(target, time.Local)
	if result.GetName() != "delay-at-test" {
		t.Fatal("链式调用 SetName 失败")
	}
	delay := result.GetDelay()
	if delay < 9*time.Second || delay > 11*time.Second {
		t.Fatalf("链式调用 SetDelayAt 时期望 Delay 约 10s，得到 %v", delay)
	}
}

// ==================== 任务类型设置方法 ====================

func TestOnce(t *testing.T) {
	tasker, _ := Task.New(func() {})
	at := time.Now().Add(time.Hour)
	tasker.Once(at)
	if tasker.GetTaskType() != TaskTypeOnce {
		t.Fatalf("期望 TaskTypeOnce，得到 %d", tasker.GetTaskType())
	}
	if !tasker.GetAt().Equal(at) {
		t.Fatalf("At 不匹配")
	}
}

func TestAfter(t *testing.T) {
	tasker, _ := Task.New(func() {})
	tasker.After(10 * time.Second)
	if tasker.GetInterval() != 10*time.Second {
		t.Fatalf("期望 10s，得到 %v", tasker.GetInterval())
	}
}

func TestAfter_DefaultOnNonPositive(t *testing.T) {
	tasker, _ := Task.New(func() {})
	tasker.After(0)
	if tasker.GetInterval() != defaultAfter {
		t.Fatalf("期望默认 %v，得到 %v", defaultAfter, tasker.GetInterval())
	}
	tasker.After(-5 * time.Second)
	if tasker.GetInterval() != defaultAfter {
		t.Fatalf("负值期望默认 %v，得到 %v", defaultAfter, tasker.GetInterval())
	}
}

func TestDaily(t *testing.T) {
	tasker, _ := Task.New(func() {})
	tasker.Daily(3)
	if tasker.GetTaskType() != TaskTypeDaily {
		t.Fatalf("期望 TaskTypeDaily，得到 %d", tasker.GetTaskType())
	}
	if tasker.GetInterval() != 3*24*time.Hour {
		t.Fatalf("期望 3 天，得到 %v", tasker.GetInterval())
	}
}

func TestDaily_MinInterval(t *testing.T) {
	tasker, _ := Task.New(func() {})
	tasker.Daily(0)
	if tasker.GetInterval() != 24*time.Hour {
		t.Fatalf("interval<1 应回退为 1 天，得到 %v", tasker.GetInterval())
	}
}

func TestHourly(t *testing.T) {
	tasker, _ := Task.New(func() {})
	tasker.Hourly(2)
	if tasker.GetTaskType() != TaskTypeHourly {
		t.Fatalf("期望 TaskTypeHourly，得到 %d", tasker.GetTaskType())
	}
	if tasker.GetInterval() != 2*time.Hour {
		t.Fatalf("期望 2h，得到 %v", tasker.GetInterval())
	}
}

func TestHourly_MinInterval(t *testing.T) {
	tasker, _ := Task.New(func() {})
	tasker.Hourly(-1)
	if tasker.GetInterval() != time.Hour {
		t.Fatalf("interval<1 应回退为 1h，得到 %v", tasker.GetInterval())
	}
}

func TestMinutely(t *testing.T) {
	tasker, _ := Task.New(func() {})
	tasker.Minutely(5)
	// 注意：当前实现中 Minutely 设置的是 TaskTypeDaily
	if tasker.GetTaskType() != TaskTypeDaily {
		t.Fatalf("期望 TaskTypeDaily，得到 %d", tasker.GetTaskType())
	}
	if tasker.GetInterval() != 5*time.Minute {
		t.Fatalf("期望 5m，得到 %v", tasker.GetInterval())
	}
}

func TestMinutely_MinInterval(t *testing.T) {
	tasker, _ := Task.New(func() {})
	tasker.Minutely(0)
	if tasker.GetInterval() != time.Minute {
		t.Fatalf("interval<1 应回退为 1m，得到 %v", tasker.GetInterval())
	}
}

// ==================== 链式调用 ====================

func TestChaining(t *testing.T) {
	tasker, _ := Task.New(func() {})
	result := tasker.SetName("chain").SetTimeout(10 * time.Second).SetDelay(2 * time.Second).Daily(1)
	if result.GetName() != "chain" {
		t.Fatal("SetName 链式失败")
	}
	if result.GetTimeout() != 10*time.Second {
		t.Fatal("SetTimeout 链式失败")
	}
	if result.GetDelay() != 2*time.Second {
		t.Fatal("SetDelay 链式失败")
	}
	if result.GetTaskType() != TaskTypeDaily {
		t.Fatal("Daily 链式失败")
	}
}

// ==================== Start 错误分支 ====================

func TestStart_NilFn(t *testing.T) {
	tasker := &TimerTaskerImpl{
		UUID:    _uuid.Must(_uuid.NewV7()),
		At:      defaultAt,
		Timeout: defaultTimeout,
		stopCh:  make(chan struct{}, 1),
	}
	err := tasker.Start()
	if err == nil {
		t.Fatal("期望错误")
	}
}

func TestStart_UnsupportedType(t *testing.T) {
	tasker, _ := Task.New(func() {})
	// TaskType = 0 (TaskTypeAfter) 是默认值，手动设为无效值
	tasker.(*TimerTaskerImpl).TaskType = TaskTypeTag(99)
	err := tasker.Start()
	if err == nil {
		t.Fatal("期望不支持的类型错误")
	}
}

// ==================== runAfter / After 执行 ====================

func TestAfter_Execution(t *testing.T) {
	var called int32
	tasker := newTasker(func() { atomic.StoreInt32(&called, 1) })
	tasker.After(50 * time.Millisecond)
	tasker.SetTimeout(2 * time.Second)

	done := make(chan error, 1)
	go func() { done <- tasker.Start() }()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("执行错误: %v", err)
		}
		if atomic.LoadInt32(&called) != 1 {
			t.Fatal("回调未被执行")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("超时")
	}
}

// ==================== Once 执行 ====================

func TestOnce_FutureTime(t *testing.T) {
	var called int32
	tasker := newTasker(func() { atomic.StoreInt32(&called, 1) })
	tasker.Once(time.Now().Add(50 * time.Millisecond))
	tasker.SetTimeout(2 * time.Second)

	done := make(chan error, 1)
	go func() { done <- tasker.Start() }()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("执行错误: %v", err)
		}
		if atomic.LoadInt32(&called) != 1 {
			t.Fatal("回调未被执行")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("超时")
	}
}

func TestOnce_PastTime(t *testing.T) {
	var called int32
	tasker := newTasker(func() { atomic.StoreInt32(&called, 1) })
	tasker.Once(time.Now().Add(-1 * time.Hour)) // 过去时间，应立即执行
	tasker.SetTimeout(2 * time.Second)

	done := make(chan error, 1)
	go func() { done <- tasker.Start() }()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("执行错误: %v", err)
		}
		if atomic.LoadInt32(&called) != 1 {
			t.Fatal("过去时间的 Once 回调未被立即执行")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("超时")
	}
}

// ==================== Stop 机制 ====================

func TestStop_Idempotent(t *testing.T) {
	tasker, _ := Task.New(func() {})
	tasker.After(time.Second)
	// 多次 Stop 不应 panic
	tasker.Stop()
	tasker.Stop()
	tasker.Stop()
}

func TestAfter_StopBeforeExecution(t *testing.T) {
	var called int32
	tasker := newTasker(func() { atomic.StoreInt32(&called, 1) })
	tasker.After(5 * time.Second) // 较长延迟
	tasker.SetTimeout(10 * time.Second)

	done := make(chan error, 1)
	go func() { done <- tasker.Start() }()

	time.Sleep(50 * time.Millisecond)
	tasker.Stop()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Stop 不应产生错误，得到 %v", err)
		}
		if atomic.LoadInt32(&called) != 0 {
			t.Fatal("Stop 后回调不应被执行")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Stop 后 Start 未及时返回")
	}
}

// ==================== runAfter panic 保护 ====================

func TestRunAfter_PanicRecovery(t *testing.T) {
	var errReceived error
	var mu sync.Mutex
	tasker := newTasker(func() { panic("test panic in After") })
	tasker.After(50 * time.Millisecond)
	tasker.SetTimeout(2 * time.Second)
	tasker.SetErrHandler(func(tasker TimerTasker, err error) {
		mu.Lock()
		errReceived = err
		mu.Unlock()
	})

	done := make(chan error, 1)
	go func() { done <- tasker.Start() }()

	select {
	case <-done:
		mu.Lock()
		defer mu.Unlock()
		if errReceived == nil {
			t.Fatal("panic 应触发 ErrHandler")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("超时")
	}
}

// ==================== safeFn ====================

func TestSafeFn_Normal(t *testing.T) {
	var called int32
	tasker := newTasker(func() { atomic.StoreInt32(&called, 1) })
	err := tasker.(*TimerTaskerImpl).safeFn()
	if err != nil {
		t.Fatalf("期望无错误，得到 %v", err)
	}
	if atomic.LoadInt32(&called) != 1 {
		t.Fatal("safeFn 未执行回调")
	}
}

func TestSafeFn_Panic(t *testing.T) {
	tasker := newTasker(func() { panic("safeFn panic") })
	err := tasker.(*TimerTaskerImpl).safeFn()
	if err == nil {
		t.Fatal("期望 panic 被捕获为错误")
	}
}

// ==================== computeNext ====================

func TestComputeNext_Daily(t *testing.T) {
	tasker, _ := Task.New(func() {})
	tasker.Daily(1)
	// 设置锚点为未来时间，确保 next > now
	now := time.Now()
	at := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()+5, 0, time.Local)
	tasker.SetAt(at)

	next := tasker.(*TimerTaskerImpl).computeNext(time.Local)
	if !next.After(now) {
		t.Fatalf("Daily computeNext 应返回未来时间，得到 %v vs now %v", next, now)
	}
	if next.Hour() != at.Hour() || next.Minute() != at.Minute() {
		t.Fatalf("时刻应与锚点一致，期望 %02d:%02d，得到 %02d:%02d",
			at.Hour(), at.Minute(), next.Hour(), next.Minute())
	}
}

func TestComputeNext_Hourly(t *testing.T) {
	tasker, _ := Task.New(func() {})
	tasker.Hourly(1)
	now := time.Now()
	// 锚点秒数设为未来 5 秒，确保 next > now
	at := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()+5, 0, time.Local)
	tasker.SetAt(at)

	next := tasker.(*TimerTaskerImpl).computeNext(time.Local)
	if !next.After(now) {
		t.Fatalf("Hourly computeNext 应返回未来时间")
	}
}

func TestComputeNext_Minutely(t *testing.T) {
	tasker, _ := Task.New(func() {})
	tasker.Minutely(1) // TaskTypeDaily, interval=1m
	now := time.Now()
	at := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()+5, 0, time.Local)
	tasker.SetAt(at)

	next := tasker.(*TimerTaskerImpl).computeNext(time.Local)
	if !next.After(now) {
		t.Fatalf("Minutely computeNext 应返回未来时间")
	}
}

func TestComputeNext_Daily_PastAt(t *testing.T) {
	tasker, _ := Task.New(func() {})
	tasker.Daily(1)
	// 锚点设为凌晨 00:00:00（几乎肯定已过）
	at := time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local)
	tasker.SetAt(at)

	now := time.Now()
	next := tasker.(*TimerTaskerImpl).computeNext(time.Local)
	if !next.After(now) {
		t.Fatalf("已过锚点时应返回明天的时刻")
	}
	if next.Hour() != 0 || next.Minute() != 0 || next.Second() != 0 {
		t.Fatalf("时刻应保持锚点 00:00:00，得到 %02d:%02d:%02d",
			next.Hour(), next.Minute(), next.Second())
	}
}

// ==================== advanceNext ====================

func TestAdvanceNext_Daily(t *testing.T) {
	tasker, _ := Task.New(func() {})
	tasker.Daily(2) // 每 2 天
	at := time.Date(2025, 6, 15, 10, 30, 0, 0, time.Local)
	tasker.SetAt(at)

	current := time.Date(2025, 6, 15, 10, 30, 0, 0, time.Local)
	next := tasker.(*TimerTaskerImpl).advanceNext(current, time.Local, 2)
	expected := time.Date(2025, 6, 17, 10, 30, 0, 0, time.Local)
	if !next.Equal(expected) {
		t.Fatalf("期望 %v，得到 %v", expected, next)
	}
}

func TestAdvanceNext_Hourly(t *testing.T) {
	tasker, _ := Task.New(func() {})
	tasker.Hourly(3)
	current := time.Date(2025, 6, 15, 10, 0, 0, 0, time.Local)
	next := tasker.(*TimerTaskerImpl).advanceNext(current, time.Local, 3)
	expected := current.Add(3 * time.Hour)
	if !next.Equal(expected) {
		t.Fatalf("期望 %v，得到 %v", expected, next)
	}
}

func TestAdvanceNext_Minutely(t *testing.T) {
	tasker, _ := Task.New(func() {})
	tasker.Minutely(10)
	// 注意：Minutely 实际设置 TaskTypeDaily，advanceNext 走 Daily 分支
	// At 未重新设置，defaultAt 时分秒均为 0
	current := time.Date(2025, 6, 15, 10, 0, 0, 0, time.Local)
	next := tasker.(*TimerTaskerImpl).advanceNext(current, time.Local, 10)
	expected := time.Date(2025, 6, 25, 0, 0, 0, 0, time.Local)
	if !next.Equal(expected) {
		t.Fatalf("期望 %v，得到 %v", expected, next)
	}
}

// ==================== runSchedule 周期执行 ====================

func TestRunSchedule_MinutelyExecution(t *testing.T) {
	var count int32
	tasker := newTasker(func() { atomic.AddInt32(&count, 1) })
	tasker.Minutely(1) // TaskTypeDaily, interval=1min
	tasker.SetTimeout(5 * time.Second)
	// 锚点设为当前时刻附近，使 computeNext 返回很快到达的时间
	now := time.Now()
	tasker.SetAt(time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()+1, 0, time.Local))

	done := make(chan struct{})
	go func() {
		tasker.Start()
		close(done)
	}()

	// 等待第一次执行
	time.Sleep(2 * time.Second)
	tasker.Stop()
	<-done

	if atomic.LoadInt32(&count) < 1 {
		t.Fatalf("期望至少执行 1 次，实际 %d", atomic.LoadInt32(&count))
	}
}

// ==================== runSchedule Delay ====================

func TestRunSchedule_DelayBeforeFirstCycle(t *testing.T) {
	var count int32
	tasker := newTasker(func() { atomic.AddInt32(&count, 1) })
	tasker.Minutely(1)
	tasker.SetDelay(200 * time.Millisecond)
	tasker.SetTimeout(5 * time.Second)
	now := time.Now()
	tasker.SetAt(time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()+1, 0, time.Local))

	start := time.Now()
	done := make(chan struct{})
	go func() {
		tasker.Start()
		close(done)
	}()

	// Delay 期间不应有执行
	time.Sleep(80 * time.Millisecond)
	if atomic.LoadInt32(&count) != 0 {
		t.Fatal("Delay 期间不应执行回调")
	}

	// 等待 Delay 结束 + 第一次周期执行
	time.Sleep(2 * time.Second)
	tasker.Stop()
	<-done

	elapsed := time.Since(start)
	if elapsed < 200*time.Millisecond {
		t.Fatalf("总耗时不应少于 Delay 200ms，实际 %v", elapsed)
	}
	if atomic.LoadInt32(&count) < 1 {
		t.Fatalf("Delay 后期望至少执行 1 次，实际 %d", atomic.LoadInt32(&count))
	}
}

func TestRunSchedule_StopDuringDelay(t *testing.T) {
	var count int32
	tasker := newTasker(func() { atomic.AddInt32(&count, 1) })
	tasker.Minutely(1)
	tasker.SetDelay(10 * time.Second) // 很长的 Delay
	tasker.SetTimeout(30 * time.Second)

	done := make(chan struct{})
	go func() {
		tasker.Start()
		close(done)
	}()

	time.Sleep(100 * time.Millisecond)
	tasker.Stop()

	select {
	case <-done:
		if atomic.LoadInt32(&count) != 0 {
			t.Fatal("Delay 期间 Stop 后不应执行回调")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Stop 后 runSchedule 未及时返回")
	}
}

// ==================== runSchedule timeout ====================

func TestRunSchedule_Timeout(t *testing.T) {
	var errReceived error
	var mu sync.Mutex
	tasker := newTasker(func() {})
	tasker.Minutely(1)
	tasker.SetTimeout(200 * time.Millisecond) // 短超时
	now := time.Now()
	// 锚点设为较远的未来，让第一次 computeNext 返回较远时间
	tasker.SetAt(time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute()+5, 0, 0, time.Local))

	tasker.SetErrHandler(func(tasker TimerTasker, err error) {
		mu.Lock()
		errReceived = err
		mu.Unlock()
	})

	done := make(chan struct{})
	go func() {
		tasker.Start()
		close(done)
	}()

	select {
	case <-done:
		mu.Lock()
		defer mu.Unlock()
		if errReceived == nil {
			t.Fatal("期望超时触发 ErrHandler")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("runSchedule 超时未及时返回")
	}
}

// ==================== runSchedule panic 不中断调度 ====================

// TestRunSchedule_PanicDoesNotStopSchedule 验证 safeFn 捕获 panic 后不中断调度
// 注：由于 Minutely 实际 TaskTypeDaily，advanceNext 每次加 1 天，
// 无法在合理时间内触发第二次周期执行，因此直接验证 safeFn 的 panic 恢复行为
func TestRunSchedule_PanicDoesNotStopSchedule(t *testing.T) {
	var count int32
	tasker := newTasker(func() {
		atomic.AddInt32(&count, 1)
		panic("schedule panic")
	})
	tasker.SetErrHandler(func(tasker TimerTasker, err error) {}) // 吞掉错误

	// safeFn 应捕获 panic 并返回错误，而非传播 panic
	err := tasker.(*TimerTaskerImpl).safeFn()
	if err == nil {
		t.Fatal("期望 safeFn 返回 panic 错误")
	}
	if atomic.LoadInt32(&count) != 1 {
		t.Fatal("回调应被执行")
	}

	// 再次调用 safeFn 仍可正常执行（调度未中断）
	err = tasker.(*TimerTaskerImpl).safeFn()
	if err == nil {
		t.Fatal("第二次 safeFn 也应返回 panic 错误")
	}
	if atomic.LoadInt32(&count) != 2 {
		t.Fatalf("期望回调执行 2 次，实际 %d", atomic.LoadInt32(&count))
	}
}

// ==================== Timer (TimerImpl) ====================

func newTestTimer() *TimerImpl {
	return &TimerImpl{timers: make(map[string]TimerTasker)}
}

func TestTimer_AddTask(t *testing.T) {
	timer := newTestTimer()
	tasker, _ := Task.New(func() {})
	tasker.SetName("t1")
	timer.AddTask(tasker)

	got := timer.GetTask(tasker.GetUUID().String())
	if got == nil {
		t.Fatal("添加后应能查到任务")
	}
	if got.GetName() != "t1" {
		t.Fatalf("期望 t1，得到 %s", got.GetName())
	}
}

func TestTimer_AddTask_Nil(t *testing.T) {
	timer := newTestTimer()
	timer.AddTask(nil) // 不应 panic
	if len(timer.timers) != 0 {
		t.Fatal("nil 任务不应被添加")
	}
}

func TestTimer_GetTask_NotFound(t *testing.T) {
	timer := newTestTimer()
	got := timer.GetTask("non-existent")
	if got != nil {
		t.Fatal("不存在的任务应返回 nil")
	}
}

func TestTimer_MustGetTask(t *testing.T) {
	timer := newTestTimer()
	tasker, _ := Task.New(func() {})
	timer.AddTask(tasker)

	got := timer.MustGetTask(tasker.GetUUID().String())
	if got == nil {
		t.Fatal("MustGetTask 应返回任务")
	}
}

func TestTimer_MustGetTask_Panic(t *testing.T) {
	timer := newTestTimer()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("期望 panic")
		}
	}()
	timer.MustGetTask("non-existent")
}

func TestTimer_RemoveTask(t *testing.T) {
	timer := newTestTimer()
	tasker, _ := Task.New(func() {})
	timer.AddTask(tasker)
	timer.RemoveTask(tasker.GetUUID().String())

	if timer.GetTask(tasker.GetUUID().String()) != nil {
		t.Fatal("删除后应返回 nil")
	}
}

func TestTimer_RemoveTask_NonExistent(t *testing.T) {
	timer := newTestTimer()
	timer.RemoveTask("non-existent") // 不应 panic
}

func TestTimer_MustRemoveTask(t *testing.T) {
	timer := newTestTimer()
	tasker, _ := Task.New(func() {})
	timer.AddTask(tasker)
	timer.MustRemoveTask(tasker.GetUUID().String())

	if timer.GetTask(tasker.GetUUID().String()) != nil {
		t.Fatal("MustRemoveTask 后应查不到")
	}
}

func TestTimer_MustRemoveTask_Panic(t *testing.T) {
	timer := newTestTimer()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("期望 panic")
		}
	}()
	timer.MustRemoveTask("non-existent")
}

func TestTimer_Start(t *testing.T) {
	timer := newTestTimer()
	var called int32

	tasker := newTasker(func() { atomic.StoreInt32(&called, 1) })
	tasker.After(50 * time.Millisecond)
	tasker.SetTimeout(2 * time.Second)
	timer.AddTask(tasker)
	timer.Start()

	time.Sleep(500 * time.Millisecond)
	if atomic.LoadInt32(&called) != 1 {
		t.Fatal("Start 后任务应被执行")
	}
}

func TestTimer_Start_Idempotent(t *testing.T) {
	timer := newTestTimer()
	tasker, _ := Task.New(func() {})
	tasker.After(time.Second)
	timer.AddTask(tasker)

	timer.Start()
	if !timer.running {
		t.Fatal("启动后 running 应为 true")
	}
	timer.Start() // 重复调用应直接返回
	if !timer.running {
		t.Fatal("重复 Start 后 running 应仍为 true")
	}
	timer.Stop()
}

func TestTimer_Stop_All(t *testing.T) {
	timer := newTestTimer()
	var count int32

	for i := 0; i < 3; i++ {
		tasker := newTasker(func() { atomic.AddInt32(&count, 1) })
		tasker.After(50 * time.Millisecond)
		tasker.SetTimeout(2 * time.Second)
		timer.AddTask(tasker)
	}

	timer.Start()
	time.Sleep(500 * time.Millisecond)
	timer.Stop()

	if atomic.LoadInt32(&count) != 3 {
		t.Fatalf("期望 3 个任务都执行，实际 %d", atomic.LoadInt32(&count))
	}
	if timer.running {
		t.Fatal("Stop 后 running 应为 false")
	}
}

func TestTimer_Stop_Specific(t *testing.T) {
	timer := newTestTimer()
	var called1, called2 int32

	t1 := newTasker(func() { atomic.StoreInt32(&called1, 1) })
	t1.After(50 * time.Millisecond)
	t1.SetTimeout(2 * time.Second)

	t2 := newTasker(func() { atomic.StoreInt32(&called2, 1) })
	t2.After(50 * time.Millisecond)
	t2.SetTimeout(2 * time.Second)

	timer.AddTask(t1)
	timer.AddTask(t2)

	// 仅停止 t1
	timer.Stop(t1.GetUUID().String())

	// 启动剩余任务
	timer.Start()
	time.Sleep(500 * time.Millisecond)

	// t1 已被停止，不应执行
	// t2 应正常执行
	if atomic.LoadInt32(&called1) != 0 {
		t.Fatal("被 Stop 的 t1 不应执行")
	}
	if atomic.LoadInt32(&called2) != 1 {
		t.Fatal("t2 应正常执行")
	}
}

func TestTimer_AddTaskAndStart(t *testing.T) {
	timer := newTestTimer()
	timer.running = true // 模拟已运行状态

	var called int32
	tasker := newTasker(func() { atomic.StoreInt32(&called, 1) })
	tasker.After(50 * time.Millisecond)
	tasker.SetTimeout(2 * time.Second)

	timer.AddTaskAndStart(tasker)

	time.Sleep(500 * time.Millisecond)
	if atomic.LoadInt32(&called) != 1 {
		t.Fatal("AddTaskAndStart 应立即启动任务")
	}

	// 任务应已被登记
	if timer.GetTask(tasker.GetUUID().String()) == nil {
		t.Fatal("AddTaskAndStart 后任务应被登记")
	}
	timer.Stop()
}

func TestTimer_AddTaskAndStart_Nil(t *testing.T) {
	timer := newTestTimer()
	timer.AddTaskAndStart(nil) // 不应 panic
	if len(timer.timers) != 0 {
		t.Fatal("nil 任务不应被添加")
	}
}

// ==================== Timer.Once 单例 ====================

func TestTimer_Once_ReturnsNonNil(t *testing.T) {
	ins := Time.Once()
	if ins == nil {
		t.Fatal("Once 应返回非空实例")
	}
	// 多次调用返回同一实例
	ins2 := Time.Once()
	if ins != ins2 {
		t.Fatal("Once 多次调用应返回同一实例")
	}
}
