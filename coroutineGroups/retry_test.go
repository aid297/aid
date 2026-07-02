package coroutineGroups

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

// ==================== SetRetry 配置测试 ====================

func TestSetRetry_Config(t *testing.T) {
	g := New[int](4, WithRetryAttempts(3),
		WithTimeout(5*time.Second),
		WithRetryInterval(100*time.Millisecond),
		WithBackoff(BackoffExponential))

	impl := g.(*CoroutineGroupImpl[int])
	if impl.retry == nil {
		t.Fatal("SetRetry后retry配置不应为nil")
	}
	if impl.retry.maxAttempts != 3 {
		t.Errorf("期望maxAttempts=3, 实际%d", impl.retry.maxAttempts)
	}
	if impl.retry.timeout != 5*time.Second {
		t.Errorf("期望timeout=5s, 实际%v", impl.retry.timeout)
	}
	if impl.retry.interval != 100*time.Millisecond {
		t.Errorf("期望interval=100ms, 实际%v", impl.retry.interval)
	}
	if impl.retry.backoff != BackoffExponential {
		t.Error("期望backoff=BackoffExponential")
	}
}

func TestSetRetry_DefaultAttempts(t *testing.T) {
	// attempts < 1 应自动设为1
	g := New[int](4, WithRetryAttempts(0))
	impl := g.(*CoroutineGroupImpl[int])
	if impl.retry.maxAttempts != 1 {
		t.Errorf("期望maxAttempts=1, 实际%d", impl.retry.maxAttempts)
	}
}

func TestSetRetry_NoRetry(t *testing.T) {
	// 未调用SetRetry时，retry应为nil，直接执行
	g := New[int](4)
	impl := g.(*CoroutineGroupImpl[int])
	if impl.retry != nil {
		t.Error("未调用SetRetry时retry应为nil")
	}
}

// ==================== calculateWait 退避策略测试 ====================

func TestCalculateWait_Fixed(t *testing.T) {
	cfg := &RetryConfig{interval: 100 * time.Millisecond, backoff: BackoffFixed}
	// 固定退避：每次都是 interval
	for attempt := 1; attempt <= 5; attempt++ {
		got := cfg.calculateWait(attempt)
		if got != 100*time.Millisecond {
			t.Errorf("attempt=%d: 期望100ms, 实际%v", attempt, got)
		}
	}
}

func TestCalculateWait_Exponential(t *testing.T) {
	cfg := &RetryConfig{interval: 100 * time.Millisecond, backoff: BackoffExponential}
	cases := []struct {
		attempt int
		expect  time.Duration
	}{
		{1, 100 * time.Millisecond}, // 100 * 2^0 = 100
		{2, 200 * time.Millisecond}, // 100 * 2^1 = 200
		{3, 400 * time.Millisecond}, // 100 * 2^2 = 400
		{4, 800 * time.Millisecond}, // 100 * 2^3 = 800
	}
	for _, c := range cases {
		got := cfg.calculateWait(c.attempt)
		if got != c.expect {
			t.Errorf("attempt=%d: 期望%v, 实际%v", c.attempt, c.expect, got)
		}
	}
}

func TestCalculateWait_ExponentialJitter(t *testing.T) {
	cfg := &RetryConfig{interval: 100 * time.Millisecond, backoff: BackoffExponentialJitter}
	// 指数退避+抖动：base=100*2^(n-1)，总等待 >= base 且 <= 2*base
	for attempt := 1; attempt <= 4; attempt++ {
		got := cfg.calculateWait(attempt)
		base := 100 * time.Millisecond * time.Duration(1<<(attempt-1))
		if got < base || got > 2*base {
			t.Errorf("attempt=%d: 期望[%v, %v], 实际%v", attempt, base, 2*base, got)
		}
	}
}

func TestCalculateWait_ZeroInterval(t *testing.T) {
	cfg := &RetryConfig{interval: 0, backoff: BackoffExponential}
	if cfg.calculateWait(1) != 0 {
		t.Error("interval=0时应返回0")
	}
}

// ==================== GO + 重试 测试 ====================

func TestGO_Retry_SuccessAfterRetries(t *testing.T) {
	var callCount int32

	results := New[int](4).SetRetry(
		WithRetryAttempts(3),
		WithRetryInterval(10*time.Millisecond),
		WithBackoff(BackoffFixed),
	).GO(
		// 前2次失败，第3次成功
		func() Result[int] {
			n := atomic.AddInt32(&callCount, 1)
			if n < 3 {
				return &ResultImpl[int]{Error: errors.New("transient error")}
			}
			return &ResultImpl[int]{Value: 42}
		},
	)

	slice := results.ToSlice()
	if len(slice) != 1 {
		t.Fatalf("期望1个结果, 实际%d", len(slice))
	}
	if !slice[0].IsOK() {
		t.Errorf("重试后应为OK, error: %v", slice[0].GetError())
	}
	if slice[0].GetValue() != 42 {
		t.Errorf("期望值42, 实际%v", slice[0].GetValue())
	}
	if atomic.LoadInt32(&callCount) != 3 {
		t.Errorf("期望调用3次, 实际%d次", atomic.LoadInt32(&callCount))
	}
}

func TestGO_Retry_AllFail(t *testing.T) {
	var callCount int32
	testErr := errors.New("always fail")

	results := New[int](4).SetRetry(
		WithRetryAttempts(3),
		WithRetryInterval(10*time.Millisecond),
		WithBackoff(BackoffFixed),
	).GO(
		func() Result[int] {
			atomic.AddInt32(&callCount, 1)
			return &ResultImpl[int]{Error: testErr}
		},
	)

	slice := results.ToSlice()
	if len(slice) != 1 {
		t.Fatalf("期望1个结果, 实际%d", len(slice))
	}
	if slice[0].IsOK() {
		t.Error("全部失败后不应为OK")
	}
	if !errors.Is(slice[0].GetError(), testErr) {
		t.Errorf("期望错误%v, 实际%v", testErr, slice[0].GetError())
	}
	if atomic.LoadInt32(&callCount) != 3 {
		t.Errorf("期望调用3次, 实际%d次", atomic.LoadInt32(&callCount))
	}
}

func TestGO_Retry_Timeout(t *testing.T) {
	var callCount int32

	results := New[int](4).SetRetry(
		WithRetryAttempts(2),
		WithTimeout(50*time.Millisecond),
		WithRetryInterval(10*time.Millisecond),
		WithBackoff(BackoffFixed),
	).GO(
		// 模拟长时间执行，每次都超时
		func() Result[int] {
			atomic.AddInt32(&callCount, 1)
			time.Sleep(200 * time.Millisecond)
			return &ResultImpl[int]{Value: 1}
		},
	)

	slice := results.ToSlice()
	if len(slice) != 1 {
		t.Fatalf("期望1个结果, 实际%d", len(slice))
	}
	if slice[0].IsOK() {
		t.Error("超时后不应为OK")
	}
	if !errors.Is(slice[0].GetError(), ErrTimeout) {
		t.Errorf("期望ErrTimeout, 实际%v", slice[0].GetError())
	}
	if atomic.LoadInt32(&callCount) != 2 {
		t.Errorf("期望调用2次(2次超时), 实际%d次", atomic.LoadInt32(&callCount))
	}
}

func TestGO_Retry_TimeoutThenSuccess(t *testing.T) {
	var callCount int32

	results := New[int](4).SetRetry(
		WithRetryAttempts(3),
		WithTimeout(50*time.Millisecond),
		WithRetryInterval(10*time.Millisecond),
		WithBackoff(BackoffFixed),
	).GO(
		// 第1次超时，第2次快速成功
		func() Result[int] {
			n := atomic.AddInt32(&callCount, 1)
			if n == 1 {
				time.Sleep(200 * time.Millisecond) // 超时
				return &ResultImpl[int]{Value: 1}
			}
			return &ResultImpl[int]{Value: 99}
		},
	)

	slice := results.ToSlice()
	if len(slice) != 1 {
		t.Fatalf("期望1个结果, 实际%d", len(slice))
	}
	if !slice[0].IsOK() {
		t.Errorf("超时后重试成功应为OK, error: %v", slice[0].GetError())
	}
	if slice[0].GetValue() != 99 {
		t.Errorf("期望值99, 实际%v", slice[0].GetValue())
	}
	if atomic.LoadInt32(&callCount) != 2 {
		t.Errorf("期望调用2次(1次超时+1次成功), 实际%d次", atomic.LoadInt32(&callCount))
	}
}

func TestGO_Retry_SkipNotRetried(t *testing.T) {
	var callCount int32

	results := New[int](4).SetRetry(
		WithRetryAttempts(3),
		WithRetryInterval(10*time.Millisecond),
		WithBackoff(BackoffFixed),
	).GO(
		// 返回Skip，不应触发重试
		func() Result[int] {
			atomic.AddInt32(&callCount, 1)
			return &ResultImpl[int]{Value: 1, Skip: true}
		},
	)

	slice := results.ToSlice()
	if len(slice) != 1 {
		t.Fatalf("期望1个结果, 实际%d", len(slice))
	}
	if !slice[0].IsSkip() {
		t.Error("应为Skip")
	}
	if atomic.LoadInt32(&callCount) != 1 {
		t.Errorf("Skip不应重试, 期望调用1次, 实际%d次", atomic.LoadInt32(&callCount))
	}
}

func TestGO_Retry_MultipleGoroutinesIndependent(t *testing.T) {
	// 每个协程独立重试，互不影响
	var callCount int32

	results := New[int](2).SetRetry(
		WithRetryAttempts(3),
		WithRetryInterval(10*time.Millisecond),
		WithBackoff(BackoffFixed),
	).GO(
		func() Result[int] {
			n := atomic.AddInt32(&callCount, 1)
			if n < 3 {
				return &ResultImpl[int]{Error: errors.New("fail")}
			}
			return &ResultImpl[int]{Value: 1}
		},
		func() Result[int] {
			return &ResultImpl[int]{Value: 2} // 一次成功
		},
	)

	slice := results.ToSlice()
	if len(slice) != 2 {
		t.Fatalf("期望2个结果, 实际%d", len(slice))
	}

	okCount := 0
	for _, r := range slice {
		if r.IsOK() {
			okCount++
		}
	}
	if okCount != 2 {
		t.Errorf("期望2个OK结果, 实际%d", okCount)
	}
}

// ==================== GO + 退避策略时间验证 ====================

func TestGO_Retry_BackoffFixedTiming(t *testing.T) {
	var callCount int32

	start := time.Now()
	New[int](4).SetRetry(
		WithRetryAttempts(3),
		WithRetryInterval(50*time.Millisecond),
		WithBackoff(BackoffFixed),
	).GO(
		func() Result[int] {
			atomic.AddInt32(&callCount, 1)
			return &ResultImpl[int]{Error: errors.New("fail")}
		},
	)
	elapsed := time.Since(start)

	// 3次执行 + 2次固定间隔50ms = 至少100ms
	if elapsed < 100*time.Millisecond {
		t.Errorf("固定退避3次应至少耗时100ms, 实际%v", elapsed)
	}
	// 上限放宽：不应超过500ms
	if elapsed > 500*time.Millisecond {
		t.Errorf("固定退避3次耗时过长: %v", elapsed)
	}
}

func TestGO_Retry_BackoffExponentialTiming(t *testing.T) {
	var callCount int32

	start := time.Now()
	New[int](4).SetRetry(
		WithRetryAttempts(4),
		WithRetryInterval(50*time.Millisecond),
		WithBackoff(BackoffExponential),
	).GO(
		func() Result[int] {
			atomic.AddInt32(&callCount, 1)
			return &ResultImpl[int]{Error: errors.New("fail")}
		},
	)
	elapsed := time.Since(start)

	// 4次执行 + 3次指数退避: 50 + 100 + 200 = 350ms
	// 至少350ms
	if elapsed < 350*time.Millisecond {
		t.Errorf("指数退避4次应至少耗时350ms, 实际%v", elapsed)
	}
	// 上限放宽
	if elapsed > 800*time.Millisecond {
		t.Errorf("指数退避4次耗时过长: %v", elapsed)
	}
}

// ==================== GOBatch + 重试 测试 ====================

func TestGOBatch_Retry_SuccessAfterRetries(t *testing.T) {
	var callCount int32

	// total=4, capacity=2 -> 2批次 x 2容量 = 4个结果
	results, err := New[int](4).SetRetry(
		WithRetryAttempts(3),
		WithRetryInterval(10*time.Millisecond),
		WithBackoff(BackoffFixed),
	).GOBatch(4, 2, func(batch, capacity uint) Result[int] {
		n := atomic.AddInt32(&callCount, 1)
		if n < 3 {
			return &ResultImpl[int]{Error: errors.New("transient")}
		}
		return &ResultImpl[int]{Value: int(batch*2 + capacity)}
	})

	if err != nil {
		t.Fatalf("不应返回错误: %v", err)
	}
	if results.Length() != 4 {
		t.Errorf("期望4个结果, 实际%d", results.Length())
	}

	okCount := 0
	for _, r := range results.ToSlice() {
		if r.IsOK() {
			okCount++
		}
	}
	if okCount != 4 {
		t.Errorf("期望4个OK结果, 实际%d", okCount)
	}
}

func TestGOBatch_Retry_Timeout(t *testing.T) {
	var callCount int32

	results, err := New[int](4).SetRetry(
		WithRetryAttempts(2),
		WithTimeout(50*time.Millisecond),
		WithRetryInterval(10*time.Millisecond),
		WithBackoff(BackoffFixed),
	).GOBatch(1, 1, func(batch, capacity uint) Result[int] {
		atomic.AddInt32(&callCount, 1)
		time.Sleep(200 * time.Millisecond)
		return &ResultImpl[int]{Value: 1}
	})

	if err != nil {
		t.Fatalf("不应返回错误: %v", err)
	}
	slice := results.ToSlice()
	if len(slice) != 1 {
		t.Fatalf("期望1个结果, 实际%d", len(slice))
	}
	if !errors.Is(slice[0].GetError(), ErrTimeout) {
		t.Errorf("期望ErrTimeout, 实际%v", slice[0].GetError())
	}
	if atomic.LoadInt32(&callCount) != 2 {
		t.Errorf("期望调用2次, 实际%d次", atomic.LoadInt32(&callCount))
	}
}

func TestGOBatch_Retry_AllFail(t *testing.T) {
	var callCount int32
	testErr := errors.New("always fail")

	results, err := New[int](4).SetRetry(
		WithRetryAttempts(3),
		WithRetryInterval(10*time.Millisecond),
		WithBackoff(BackoffFixed),
	).GOBatch(2, 1, func(batch, capacity uint) Result[int] {
		atomic.AddInt32(&callCount, 1)
		return &ResultImpl[int]{Error: testErr}
	})

	if err != nil {
		t.Fatalf("不应返回错误: %v", err)
	}
	if results.Length() != 2 {
		t.Errorf("期望2个结果, 实际%d", results.Length())
	}
	for _, r := range results.ToSlice() {
		if r.IsOK() {
			t.Error("全部失败后不应为OK")
		}
	}
	// 2个任务各重试3次 = 6次
	if atomic.LoadInt32(&callCount) != 6 {
		t.Errorf("期望调用6次(2任务x3次), 实际%d次", atomic.LoadInt32(&callCount))
	}
}

// ==================== 向后兼容测试 ====================

func TestGO_NoRetry_BackwardCompatible(t *testing.T) {
	// 不调用SetRetry时，GO行为与之前完全一致
	results := New[int](4).GO(
		func() Result[int] { return &ResultImpl[int]{Value: 1} },
		func() Result[int] { return &ResultImpl[int]{Value: 2} },
	)

	if results.Length() != 2 {
		t.Errorf("期望2个结果, 实际%d", results.Length())
	}
}

func TestGOBatch_NoRetry_BackwardCompatible(t *testing.T) {
	results, err := New[int](4).GOBatch(4, 2, func(batch, capacity uint) Result[int] {
		return &ResultImpl[int]{Value: int(batch*2 + capacity)}
	})

	if err != nil {
		t.Fatalf("不应返回错误: %v", err)
	}
	if results.Length() != 4 {
		t.Errorf("期望4个结果, 实际%d", results.Length())
	}
}
