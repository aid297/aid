package coroutineGroupV2

import (
	"errors"
	"sync/atomic"
	"testing"
)

// ==================== ResultImpl 方法测试 ====================

func TestResultImpl_IsOK(t *testing.T) {
	// 正常结果
	r := &ResultImpl[int]{Value: 42}
	if !r.IsOK() {
		t.Error("无错误且未跳过时应返回true")
	}

	// 有错误
	r.Error = errors.New("err")
	if r.IsOK() {
		t.Error("有错误时应返回false")
	}

	// 无错误但跳过
	r.Error = nil
	r.Skip = true
	if r.IsOK() {
		t.Error("跳过时应返回false")
	}
}

func TestResultImpl_IsSkip(t *testing.T) {
	r := &ResultImpl[int]{Skip: true}
	if !r.IsSkip() {
		t.Error("Skip为true时IsSkip应返回true")
	}
	r.Skip = false
	if r.IsSkip() {
		t.Error("Skip为false时IsSkip应返回false")
	}
}

func TestResultImpl_GetError(t *testing.T) {
	err := errors.New("test error")
	r := &ResultImpl[int]{Error: err}
	if !errors.Is(err, r.GetError()) {
		t.Error("GetError应返回设置的error")
	}

	r2 := &ResultImpl[int]{}
	if r2.GetError() != nil {
		t.Error("未设置error时应返回nil")
	}
}

func TestResultImpl_GetValue(t *testing.T) {
	r := &ResultImpl[string]{Value: "hello"}
	if r.GetValue() != "hello" {
		t.Errorf("期望hello, 实际%v", r.GetValue())
	}

	r2 := &ResultImpl[int]{Value: 99}
	if r2.GetValue() != 99 {
		t.Errorf("期望99, 实际%v", r2.GetValue())
	}
}

func TestResultImpl_SetOK(t *testing.T) {
	r := &ResultImpl[int]{Error: errors.New("err")}
	r.SetOK(true)
	if r.GetError() != nil {
		t.Error("SetOK后Error应为nil")
	}
}

func TestResultImpl_SetSkip(t *testing.T) {
	r := &ResultImpl[int]{}
	r.SetSkip(true)
	if !r.IsSkip() {
		t.Error("SetSkip(true)后IsSkip应返回true")
	}
	r.SetSkip(false)
	if r.IsSkip() {
		t.Error("SetSkip(false)后IsSkip应返回false")
	}
}

// ==================== New 构造函数测试 ====================

func TestNew(t *testing.T) {
	// 自定义并发限制
	g := New[int](4)
	if g == nil {
		t.Fatal("New应返回非nil")
	}

	// limit为0时默认设为4
	g2 := New[int](0)
	if g2 == nil {
		t.Fatal("New(0)应返回非nil")
	}
}

func TestCoroutineGroupImpl_New(t *testing.T) {
	g := New[int](4)
	g2 := g.New(8)
	if g2 == nil {
		t.Fatal("实例方法New应返回非nil")
	}
}

// ==================== GO 方法测试 ====================

func TestGO_Basic(t *testing.T) {
	results := New[int](4).
		GO(
			func() Result[int] { return &ResultImpl[int]{Value: 1} },
			func() Result[int] { return &ResultImpl[int]{Value: 2} },
			func() Result[int] { return &ResultImpl[int]{Value: 3} },
			func() Result[int] { return &ResultImpl[int]{Value: 4} },
			func() Result[int] { return &ResultImpl[int]{Value: 5} },
			func() Result[int] { return &ResultImpl[int]{Value: 6} },
			func() Result[int] { return &ResultImpl[int]{Value: 7} },
			func() Result[int] { return &ResultImpl[int]{Value: 8} },
			func() Result[int] { return &ResultImpl[int]{Value: 9} },
		)

	if results.Length() != 9 {
		t.Errorf("期望9个结果, 实际%d", results.Length())
	}

	// 验证所有结果都成功
	for _, result := range results.ToSlice() {
		if !result.IsOK() {
			t.Error("每个结果应为OK")
		}
	}
}

func TestGO_EmptyFuncs(t *testing.T) {
	results := New[int](4).GO()
	if results.Length() != 0 {
		t.Errorf("无函数时应返回空结果, 实际长度%d", results.Length())
	}
}

func TestGO_WithError(t *testing.T) {
	testErr := errors.New("test error")
	results := New[string](2).GO(
		func() Result[string] { return &ResultImpl[string]{Value: "ok"} },
		func() Result[string] { return &ResultImpl[string]{Error: testErr} },
	)

	// 并发执行结果顺序不确定，按统计数量断言
	slice := results.ToSlice()
	if len(slice) != 2 {
		t.Fatalf("期望2个结果, 实际%d", len(slice))
	}

	okCount := 0
	errCount := 0
	for _, r := range slice {
		if r.IsOK() {
			okCount++
		} else {
			errCount++
		}
	}
	if okCount != 1 {
		t.Errorf("期望1个OK结果, 实际%d", okCount)
	}
	if errCount != 1 {
		t.Errorf("期望1个错误结果, 实际%d", errCount)
	}
}

func TestGO_WithSkip(t *testing.T) {
	results := New[int](2).GO(
		func() Result[int] { return &ResultImpl[int]{Value: 1} },
		func() Result[int] { return &ResultImpl[int]{Value: 2, Skip: true} },
		func() Result[int] { return &ResultImpl[int]{Value: 3} },
	)

	// 并发执行结果顺序不确定，按统计数量断言
	slice := results.ToSlice()
	if len(slice) != 3 {
		t.Fatalf("期望3个结果, 实际%d", len(slice))
	}

	skipCount := 0
	notSkipCount := 0
	for _, r := range slice {
		if r.IsSkip() {
			skipCount++
		} else {
			notSkipCount++
		}
	}
	if skipCount != 1 {
		t.Errorf("期望1个跳过结果, 实际%d", skipCount)
	}
	if notSkipCount != 2 {
		t.Errorf("期望2个未跳过结果, 实际%d", notSkipCount)
	}
}

func TestGO_ConcurrencyLimit(t *testing.T) {
	limit := 2
	var current int32
	var maxSeen int32

	funcs := make([]Func[int], 8)
	for i := 0; i < 8; i++ {
		funcs[i] = func() Result[int] {
			c := atomic.AddInt32(&current, 1)
			// 记录最大并发数
			for {
				old := atomic.LoadInt32(&maxSeen)
				if c <= old || atomic.CompareAndSwapInt32(&maxSeen, old, c) {
					break
				}
			}
			// 短暂模拟工作
			for j := 0; j < 100000; j++ {
			}
			atomic.AddInt32(&current, -1)
			return &ResultImpl[int]{Value: 1}
		}
	}

	results := New[int](uint16(limit)).GO(funcs...)

	seen := atomic.LoadInt32(&maxSeen)
	if int(seen) > limit {
		t.Errorf("最大并发数%d超过限制%d", seen, limit)
	}
	if results.Length() != 8 {
		t.Errorf("期望8个结果, 实际%d", results.Length())
	}
}

func TestGO_SetFunc(t *testing.T) {
	g := New[int](4).SetFunc(
		func() Result[int] { return &ResultImpl[int]{Value: 10} },
		func() Result[int] { return &ResultImpl[int]{Value: 20} },
	)

	// 通过GO传入新函数（覆盖SetFunc）
	results := g.GO(
		func() Result[int] { return &ResultImpl[int]{Value: 100} },
		func() Result[int] { return &ResultImpl[int]{Value: 200} },
	)

	if results.Length() != 2 {
		t.Errorf("期望2个结果, 实际%d", results.Length())
	}
}

func TestGO_SingleFunc(t *testing.T) {
	results := New[int](4).GO(
		func() Result[int] { return &ResultImpl[int]{Value: 42} },
	)

	if results.Length() != 1 {
		t.Errorf("期望1个结果, 实际%d", results.Length())
	}
	if results.ToSlice()[0].GetValue() != 42 {
		t.Errorf("期望值42, 实际%v", results.ToSlice()[0].GetValue())
	}
}

// ==================== GOBatch 方法测试 ====================

func TestGOBatch_Basic(t *testing.T) {
	// total=6, capacity=3 -> 2批次，每批3个
	results, err := New[int](4).GOBatch(6, 3, func(batch, capacity uint) Result[int] {
		return &ResultImpl[int]{Value: int(batch*3 + capacity)}
	})

	if err != nil {
		t.Fatalf("不应返回错误: %v", err)
	}
	if results.Length() != 6 {
		t.Errorf("期望6个结果, 实际%d", results.Length())
	}

	// 所有结果应为OK
	for _, r := range results.ToSlice() {
		if !r.IsOK() {
			t.Error("每个结果应为OK")
		}
	}
}

func TestGOBatch_TotalZero(t *testing.T) {
	_, err := New[int](4).GOBatch(0, 5, func(batch, capacity uint) Result[int] {
		return &ResultImpl[int]{Value: 1}
	})
	if !errors.Is(err, ErrEmptyFuncs) {
		t.Errorf("total=0时应返回ErrEmptyFuncs, 实际: %v", err)
	}
}

func TestGOBatch_CapacityZero(t *testing.T) {
	// 源码在计算 batches 时 (total+capacities-1)/capacities 会触发除零 panic
	// capacities==0 的检查在除法之后，无法到达
	defer func() {
		if r := recover(); r == nil {
			t.Error("capacities=0时应触发panic")
		}
	}()

	_, _ = New[int](4).GOBatch(10, 0, func(batch, capacity uint) Result[int] {
		return &ResultImpl[int]{Value: 1}
	})
}

func TestGOBatch_CapacityGreaterThanTotal(t *testing.T) {
	// total=3, capacity=10 -> 1批次，3个元素（但循环按capacity=10跑，fn会被调用10次）
	results, err := New[int](4).GOBatch(3, 10, func(batch, capacity uint) Result[int] {
		return &ResultImpl[int]{Value: int(batch*uint(10) + capacity)}
	})

	if err != nil {
		t.Fatalf("不应返回错误: %v", err)
	}
	// batches=1, capacity循环10次 -> 10个结果
	if results.Length() != 10 {
		t.Errorf("期望10个结果, 实际%d", results.Length())
	}
}

func TestGOBatch_WithError(t *testing.T) {
	results, err := New[int](4).GOBatch(4, 2, func(batch, capacity uint) Result[int] {
		if batch == 0 && capacity == 0 {
			return &ResultImpl[int]{Error: errors.New("batch 0 cap 0 error")}
		}
		return &ResultImpl[int]{Value: int(batch*2 + capacity)}
	})

	if err != nil {
		t.Fatalf("GOBatch不应返回错误: %v", err)
	}

	hasErr := false
	for _, r := range results.ToSlice() {
		if !r.IsOK() {
			hasErr = true
			break
		}
	}
	if !hasErr {
		t.Error("应有至少一个错误结果")
	}
}

func TestGOBatch_SingleBatch(t *testing.T) {
	// total=3, capacity=3 -> 1批次
	results, err := New[int](2).GOBatch(3, 3, func(batch, capacity uint) Result[int] {
		return &ResultImpl[int]{Value: int(capacity)}
	})

	if err != nil {
		t.Fatalf("不应返回错误: %v", err)
	}
	if results.Length() != 3 {
		t.Errorf("期望3个结果, 实际%d", results.Length())
	}
}

func TestGOBatch_MultipleBatches(t *testing.T) {
	// total=7, capacity=3 -> 3批次（ceil(7/3)=3），每批循环3次 -> 9个结果
	var count int32
	results, err := New[int](4).GOBatch(7, 3, func(batch, capacity uint) Result[int] {
		atomic.AddInt32(&count, 1)
		return &ResultImpl[int]{Value: int(batch*3 + capacity)}
	})

	if err != nil {
		t.Fatalf("不应返回错误: %v", err)
	}
	if results.Length() != 9 {
		t.Errorf("期望9个结果, 实际%d", results.Length())
	}
}

func TestGOBatch_TotalOne(t *testing.T) {
	// total=1, capacity=1 -> 1批次
	results, err := New[int](2).GOBatch(1, 1, func(batch, capacity uint) Result[int] {
		return &ResultImpl[int]{Value: 1}
	})

	if err != nil {
		t.Fatalf("不应返回错误: %v", err)
	}
	if results.Length() != 1 {
		t.Errorf("期望1个结果, 实际%d", results.Length())
	}
}
