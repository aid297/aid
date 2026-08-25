package timers_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	_time "time"

	_uuid "github.com/google/uuid"

	_timers "github.com/aid297/aid/v2/timers"
)

func newTimer() _timers.Timer { return _timers.Time.Once().Start() }

type es struct{}

func Test1(t *testing.T) {
	closeCh := make(chan es, 1)
	timer := newTimer()
	defer timer.Stop()

	loc, _ := _time.LoadLocation("Asia/Shanghai")
	startTime := _time.Date(2026, 8, 25, 13, 47, 30, 0, loc)

	tasker := _timers.Task.New(func() { t.Log("执行定时器"); closeCh <- es{} }).Once(startTime)

	timer.AddTaskAndStart(tasker)

	ctx, cancel := context.WithTimeout(context.Background(), 10*_time.Minute)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Log("测试超时")
		return
	case <-closeCh:
		t.Log("测试成功")
	}
}

func Test2(t *testing.T) {
	closeCh := make(chan es, 1)
	timer := newTimer()
	defer timer.Stop()

	loc, _ := _time.LoadLocation("Asia/Shanghai")
	startTime := time.Now().In(loc).Add(5 * time.Second)

	tasker := _timers.Task.New(func() { t.Log("执行定时器"); closeCh <- es{} }).SetDelayAt(startTime).Minutely(1)
	timer.AddTaskAndStart(tasker)

	ctx, cancel := context.WithTimeout(context.Background(), 2*_time.Minute)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Log("测试超时")
		return
	case <-closeCh:
		t.Log("测试成功")
	}
}

// ==================== TimerTasker 构造与链式方法 ====================

func TestNew_ValidFn(t *testing.T) {
	tasker := _timers.Task.New(func() {})
	if tasker.GetUUID() == _uuid.Nil {
		t.Fatal("期望非零 UUID")
	}
	if tasker.GetTimeout() != 24*time.Hour {
		t.Fatalf("期望默认超时 24h，得到 %v", tasker.GetTimeout())
	}
	expectedAt := time.Date(0, 0, 0, 0, 0, 0, 0, time.Local)
	if tasker.GetAt() != expectedAt {
		t.Fatalf("期望默认 At %v，得到 %v", expectedAt, tasker.GetAt())
	}
}

func TestNew_NilFn(t *testing.T) {
	tasker := _timers.Task.New(nil)
	// New(nil) 不再返回 error，而是设置空函数
	if tasker == nil {
		t.Fatal("期望返回非空 tasker")
	}
	if tasker.GetFn() == nil {
		t.Fatal("nil fn 应被替换为空函数")
	}
}

func TestSetTimeout(t *testing.T) {
	tasker := _timers.Task.New(func() {})
	tasker.SetTimeout(5 * time.Second)
	if tasker.GetTimeout() != 5*time.Second {
		t.Fatalf("期望 5s，得到 %v", tasker.GetTimeout())
	}
	// 零值回退默认
	tasker.SetTimeout(0)
	if tasker.GetTimeout() != 24*time.Hour {
		t.Fatalf("期望默认超时 24h，得到 %v", tasker.GetTimeout())
	}
	// 负值回退默认
	tasker.SetTimeout(-1 * time.Second)
	if tasker.GetTimeout() != 24*time.Hour {
		t.Fatalf("期望默认超时 24h，得到 %v", tasker.GetTimeout())
	}
}

func TestSetErrHandler(t *testing.T) {
	handler := _timers.ErrHandler(func(tasker _timers.TimerTasker, err error) {})
	tasker := _timers.Task.New(func() {})
	tasker.SetErrHandler(handler)
	if tasker.GetErrHandler() == nil {
		t.Fatal("期望 ErrHandler 不为空")
	}
}

func TestSetName(t *testing.T) {
	tasker := _timers.Task.New(func() {})
	tasker.SetName("test-task")
	if tasker.GetName() != "test-task" {
		t.Fatalf("期望 test-task，得到 %s", tasker.GetName())
	}
}

func TestSetDelay(t *testing.T) {
	tasker := _timers.Task.New(func() {})
	tasker.SetDelay(5 * time.Second)
	if tasker.GetDelay() != 5*time.Second {
		t.Fatalf("期望 5s，得到 %v", tasker.GetDelay())
	}
	tasker.SetDelay(0)
	if tasker.GetDelay() != 5*time.Second {
		t.Fatalf("零值不应改变 Delay，得到 %v", tasker.GetDelay())
	}
	tasker.SetDelay(-1 * time.Second)
	if tasker.GetDelay() != 5*time.Second {
		t.Fatalf("负值不应改变 Delay，得到 %v", tasker.GetDelay())
	}
}

func TestSetDelayAt_FutureTime(t *testing.T) {
	tasker := _timers.Task.New(func() {})
	target := time.Now().Add(5 * time.Second)
	tasker.SetDelayAt(target)
	if !tasker.GetAt().Equal(target) {
		t.Fatalf("期望 At = %v，得到 %v", target, tasker.GetAt())
	}
}

func TestSetDelayAt_PastTime(t *testing.T) {
	tasker := _timers.Task.New(func() {})
	target := time.Now().Add(-1 * time.Hour)
	tasker.SetDelayAt(target)
	if !tasker.GetAt().Equal(target) {
		t.Fatalf("过去时间也应设置 At，期望 %v，得到 %v", target, tasker.GetAt())
	}
}

func TestSetDelayAt_Chaining(t *testing.T) {
	tasker := _timers.Task.New(func() {})
	target := time.Now().Add(10 * time.Second)
	result := tasker.SetName("delay-at-test").SetDelayAt(target)
	if result.GetName() != "delay-at-test" {
		t.Fatal("链式调用 SetName 失败")
	}
	if !result.GetAt().Equal(target) {
		t.Fatalf("链式调用 SetDelayAt 时期望 At = %v，得到 %v", target, result.GetAt())
	}
}

// ==================== 任务类型设置方法 ====================

func TestOnce(t *testing.T) {
	tasker := _timers.Task.New(func() {})
	at := time.Now().Add(time.Hour)
	tasker.Once(at)
	if tasker.GetTaskType() != _timers.TaskTypeOnce {
		t.Fatalf("期望 TaskTypeOnce，得到 %d", tasker.GetTaskType())
	}
	if !tasker.GetAt().Equal(at) {
		t.Fatalf("At 不匹配")
	}
}

func TestAfter(t *testing.T) {
	tasker := _timers.Task.New(func() {})
	tasker.After(10 * time.Second)
	if tasker.GetInterval() != 10*time.Second {
		t.Fatalf("期望 10s，得到 %v", tasker.GetInterval())
	}
}

func TestAfter_DefaultOnNonPositive(t *testing.T) {
	tasker := _timers.Task.New(func() {})
	tasker.After(0)
	if tasker.GetInterval() != 30*time.Second {
		t.Fatalf("期望默认 30s，得到 %v", tasker.GetInterval())
	}
	tasker.After(-5 * time.Second)
	if tasker.GetInterval() != 30*time.Second {
		t.Fatalf("负值期望默认 30s，得到 %v", tasker.GetInterval())
	}
}

func TestDaily(t *testing.T) {
	tasker := _timers.Task.New(func() {})
	tasker.Daily(3)
	if tasker.GetTaskType() != _timers.TaskTypeDaily {
		t.Fatalf("期望 TaskTypeDaily，得到 %d", tasker.GetTaskType())
	}
	if tasker.GetInterval() != 3*24*time.Hour {
		t.Fatalf("期望 3 天，得到 %v", tasker.GetInterval())
	}
}

func TestDaily_MinInterval(t *testing.T) {
	tasker := _timers.Task.New(func() {})
	tasker.Daily(0)
	if tasker.GetInterval() != 24*time.Hour {
		t.Fatalf("interval<1 应回退为 1 天，得到 %v", tasker.GetInterval())
	}
}

func TestHourly(t *testing.T) {
	tasker := _timers.Task.New(func() {})
	tasker.Hourly(2)
	if tasker.GetTaskType() != _timers.TaskTypeHourly {
		t.Fatalf("期望 TaskTypeHourly，得到 %d", tasker.GetTaskType())
	}
	if tasker.GetInterval() != 2*time.Hour {
		t.Fatalf("期望 2h，得到 %v", tasker.GetInterval())
	}
}

func TestHourly_MinInterval(t *testing.T) {
	tasker := _timers.Task.New(func() {})
	tasker.Hourly(-1)
	if tasker.GetInterval() != time.Hour {
		t.Fatalf("interval<1 应回退为 1h，得到 %v", tasker.GetInterval())
	}
}

func TestMinutely(t *testing.T) {
	tasker := _timers.Task.New(func() {})
	tasker.Minutely(5)
	if tasker.GetTaskType() != _timers.TaskTypeDaily {
		t.Fatalf("期望 TaskTypeDaily，得到 %d", tasker.GetTaskType())
	}
	if tasker.GetInterval() != 5*time.Minute {
		t.Fatalf("期望 5m，得到 %v", tasker.GetInterval())
	}
}

func TestMinutely_MinInterval(t *testing.T) {
	tasker := _timers.Task.New(func() {})
	tasker.Minutely(0)
	if tasker.GetInterval() != time.Minute {
		t.Fatalf("interval<1 应回退为 1m，得到 %v", tasker.GetInterval())
	}
}

// ==================== 链式调用 ====================

func TestChaining(t *testing.T) {
	tasker := _timers.Task.New(func() {})
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
	if result.GetTaskType() != _timers.TaskTypeDaily {
		t.Fatal("Daily 链式失败")
	}
}

// ==================== Start 错误分支 ====================

func TestStart_NilFn(t *testing.T) {
	// New(nil) 现在设置空函数，Start 会正常执行而非返回错误
	tasker := _timers.Task.New(nil)
	tasker.After(50 * time.Millisecond)
	tasker.SetTimeout(2 * time.Second)
	err := tasker.Start()
	if err != nil {
		t.Fatalf("空函数应正常执行，得到错误: %v", err)
	}
}

func TestStart_UnsupportedType(t *testing.T) {
	tasker := _timers.Task.New(func() {})
	tasker.(*_timers.TimerTaskerImpl).TaskType = _timers.TaskTypeTag(99)
	err := tasker.Start()
	if err == nil {
		t.Fatal("期望不支持的类型错误")
	}
}

// ==================== After / Once 执行 ====================

func TestAfter_Execution(t *testing.T) {
	var called int32
	tasker := _timers.Task.New(func() { atomic.StoreInt32(&called, 1) })
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

func TestOnce_FutureTime(t *testing.T) {
	var called int32
	tasker := _timers.Task.New(func() { atomic.StoreInt32(&called, 1) })
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
	tasker := _timers.Task.New(func() { atomic.StoreInt32(&called, 1) })
	tasker.Once(time.Now().Add(-1 * time.Hour))
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
	tasker := _timers.Task.New(func() {})
	tasker.After(time.Second)
	tasker.Stop()
	tasker.Stop()
	tasker.Stop()
}

func TestAfter_StopBeforeExecution(t *testing.T) {
	var called int32
	tasker := _timers.Task.New(func() { atomic.StoreInt32(&called, 1) })
	tasker.After(5 * time.Second)
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
	tasker := _timers.Task.New(func() { panic("test panic in After") })
	tasker.After(50 * time.Millisecond)
	tasker.SetTimeout(2 * time.Second)
	tasker.SetErrHandler(func(tasker _timers.TimerTasker, err error) {
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

// ==================== runSchedule 周期执行 ====================

func TestRunSchedule_MinutelyExecution(t *testing.T) {
	var count int32
	tasker := _timers.Task.New(func() { atomic.AddInt32(&count, 1) })
	tasker.Minutely(1)
	tasker.SetTimeout(5 * time.Second)
	tasker.SetDelayAt(time.Now().Add(1 * time.Second))

	done := make(chan struct{})
	go func() {
		tasker.Start()
		close(done)
	}()

	time.Sleep(2 * time.Second)
	tasker.Stop()
	<-done

	if atomic.LoadInt32(&count) < 1 {
		t.Fatalf("期望至少执行 1 次，实际 %d", atomic.LoadInt32(&count))
	}
}

// ==================== runSchedule Stop/Timeout ====================

func TestRunSchedule_StopBeforeFirstExecution(t *testing.T) {
	var count int32
	tasker := _timers.Task.New(func() { atomic.AddInt32(&count, 1) })
	tasker.Minutely(1)
	tasker.SetDelayAt(time.Now().Add(10 * time.Second))
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
			t.Fatal("Stop 后不应执行回调")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Stop 后 runSchedule 未及时返回")
	}
}

// ==================== runSchedule timeout ====================

func TestRunSchedule_Timeout(t *testing.T) {
	var errReceived error
	var mu sync.Mutex
	tasker := _timers.Task.New(func() {})
	tasker.Minutely(1)
	tasker.SetTimeout(200 * time.Millisecond)
	tasker.SetDelayAt(time.Now().Add(5 * time.Minute))

	tasker.SetErrHandler(func(tasker _timers.TimerTasker, err error) {
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

// ==================== safeFn panic 恢复（通过公开接口间接验证） ====================

func TestSafeFn_PanicRecovery(t *testing.T) {
	var count int32
	tasker := _timers.Task.New(func() {
		atomic.AddInt32(&count, 1)
		panic("schedule panic")
	})
	tasker.SetErrHandler(func(tasker _timers.TimerTasker, err error) {})

	// 通过 Start 执行一次 After 类型任务，验证 panic 被捕获后通过 ErrHandler 上报
	tasker.After(50 * time.Millisecond)
	tasker.SetTimeout(2 * time.Second)

	var errReceived error
	var mu sync.Mutex
	tasker.SetErrHandler(func(tasker _timers.TimerTasker, err error) {
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
		if atomic.LoadInt32(&count) != 1 {
			t.Fatal("回调应被执行")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("超时")
	}
}

// ==================== Timer (TimerImpl) - 通过单例测试 ====================

func TestTimer_Once_ReturnsNonNil(t *testing.T) {
	ins := _timers.Time.Once()
	if ins == nil {
		t.Fatal("Once 应返回非空实例")
	}
	ins2 := _timers.Time.Once()
	if ins != ins2 {
		t.Fatal("Once 多次调用应返回同一实例")
	}
}

func TestTimer_AddAndGetTask(t *testing.T) {
	timer := _timers.Time.Once()
	tasker := _timers.Task.New(func() {})
	tasker.SetName("get-test")
	timer.AddTask(tasker)

	got := timer.GetTask(tasker.GetUUID().String())
	if got == nil {
		t.Fatal("添加后应能查到任务")
	}
	if got.GetName() != "get-test" {
		t.Fatalf("期望 get-test，得到 %s", got.GetName())
	}

	// 清理
	timer.RemoveTask(tasker.GetUUID().String())
}

func TestTimer_AddTask_Nil(t *testing.T) {
	timer := _timers.Time.Once()
	timer.AddTask(nil) // 不应 panic
}

func TestTimer_GetTask_NotFound(t *testing.T) {
	timer := _timers.Time.Once()
	got := timer.GetTask("non-existent")
	if got != nil {
		t.Fatal("不存在的任务应返回 nil")
	}
}

func TestTimer_MustGetTask(t *testing.T) {
	timer := _timers.Time.Once()
	tasker := _timers.Task.New(func() {})
	timer.AddTask(tasker)

	got := timer.MustGetTask(tasker.GetUUID().String())
	if got == nil {
		t.Fatal("MustGetTask 应返回任务")
	}

	timer.RemoveTask(tasker.GetUUID().String())
}

func TestTimer_MustGetTask_Panic(t *testing.T) {
	timer := _timers.Time.Once()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("期望 panic")
		}
	}()
	timer.MustGetTask("non-existent")
}

func TestTimer_RemoveTask(t *testing.T) {
	timer := _timers.Time.Once()
	tasker := _timers.Task.New(func() {})
	timer.AddTask(tasker)
	timer.RemoveTask(tasker.GetUUID().String())

	if timer.GetTask(tasker.GetUUID().String()) != nil {
		t.Fatal("删除后应返回 nil")
	}
}

func TestTimer_RemoveTask_NonExistent(t *testing.T) {
	timer := _timers.Time.Once()
	timer.RemoveTask("non-existent") // 不应 panic
}

func TestTimer_MustRemoveTask(t *testing.T) {
	timer := _timers.Time.Once()
	tasker := _timers.Task.New(func() {})
	timer.AddTask(tasker)
	timer.MustRemoveTask(tasker.GetUUID().String())

	if timer.GetTask(tasker.GetUUID().String()) != nil {
		t.Fatal("MustRemoveTask 后应查不到")
	}
}

func TestTimer_MustRemoveTask_Panic(t *testing.T) {
	timer := _timers.Time.Once()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("期望 panic")
		}
	}()
	timer.MustRemoveTask("non-existent-uuid")
}

func TestTimer_StartAndStop(t *testing.T) {
	timer := _timers.Time.Once()
	timer.Stop() // 清理单例残留状态
	var called int32

	tasker := _timers.Task.New(func() { atomic.StoreInt32(&called, 1) })
	tasker.After(50 * time.Millisecond)
	tasker.SetTimeout(2 * time.Second)
	timer.AddTaskAndStart(tasker)

	time.Sleep(500 * time.Millisecond)
	if atomic.LoadInt32(&called) != 1 {
		t.Fatal("Start 后任务应被执行")
	}

	timer.Stop(tasker.GetUUID().String())
}

func TestTimer_StopAll(t *testing.T) {
	timer := _timers.Time.Once()
	timer.Stop() // 清理单例残留状态
	var count int32

	for i := 0; i < 3; i++ {
		tasker := _timers.Task.New(func() { atomic.AddInt32(&count, 1) })
		tasker.After(50 * time.Millisecond)
		tasker.SetTimeout(2 * time.Second)
		timer.AddTaskAndStart(tasker)
	}

	time.Sleep(500 * time.Millisecond)
	timer.Stop()

	if atomic.LoadInt32(&count) < 1 {
		t.Fatalf("期望至少有任务执行，实际 %d", atomic.LoadInt32(&count))
	}
}

func TestTimer_AddTaskAndStart(t *testing.T) {
	timer := _timers.Time.Once()
	var called int32

	tasker := _timers.Task.New(func() { atomic.StoreInt32(&called, 1) })
	tasker.After(50 * time.Millisecond)
	tasker.SetTimeout(2 * time.Second)
	timer.AddTaskAndStart(tasker)

	time.Sleep(500 * time.Millisecond)
	if atomic.LoadInt32(&called) != 1 {
		t.Fatal("AddTaskAndStart 应立即启动任务")
	}

	if timer.GetTask(tasker.GetUUID().String()) == nil {
		t.Fatal("AddTaskAndStart 后任务应被登记")
	}
	timer.RemoveTask(tasker.GetUUID().String())
}

func TestTimer_AddTaskAndStart_Nil(t *testing.T) {
	timer := _timers.Time.Once()
	timer.AddTaskAndStart(nil) // 不应 panic
}
