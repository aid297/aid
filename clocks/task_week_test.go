package clocks_test

import (
	"sync/atomic"
	_testing "testing"
	_time "time"

	_clocks "github.com/aid297/aid/v2/clocks"
)

// Test_WeekOfMonthOf 验证 WeekOfMonthOf 工具函数的周序号计算
func Test_WeekOfMonthOf(t *_testing.T) {
	loc := _time.UTC
	cases := []struct {
		date _time.Time
		want int
	}{
		{_time.Date(2026, 8, 1, 0, 0, 0, 0, loc), 1},  // 第1天 → 第1周
		{_time.Date(2026, 8, 7, 0, 0, 0, 0, loc), 1},  // 第7天 → 第1周
		{_time.Date(2026, 8, 8, 0, 0, 0, 0, loc), 2},  // 第8天 → 第2周
		{_time.Date(2026, 8, 14, 0, 0, 0, 0, loc), 2}, // 第14天 → 第2周
		{_time.Date(2026, 8, 15, 0, 0, 0, 0, loc), 3}, // 第15天 → 第3周
		{_time.Date(2026, 8, 21, 0, 0, 0, 0, loc), 3}, // 第21天 → 第3周
		{_time.Date(2026, 8, 22, 0, 0, 0, 0, loc), 4}, // 第22天 → 第4周
		{_time.Date(2026, 8, 28, 0, 0, 0, 0, loc), 4}, // 第28天 → 第4周
		{_time.Date(2026, 8, 29, 0, 0, 0, 0, loc), 5}, // 第29天 → 第5周
		{_time.Date(2026, 8, 31, 0, 0, 0, 0, loc), 5}, // 第31天 → 第5周
	}

	for _, c := range cases {
		got := _clocks.WeekOfMonthOf(c.date)
		if got != c.want {
			t.Errorf("WeekOfMonthOf(%s) = %d, want %d", c.date.Format("2006-01-02"), got, c.want)
		}
	}
}

// Test_WeekRuleMatch 验证规则匹配逻辑：添加当前时刻对应的规则，确认能被命中触发
func Test_WeekRuleMatch(t *_testing.T) {
	loc, _ := _time.LoadLocation("Asia/Shanghai")
	now := _time.Now().In(loc)

	// 计算当前时刻是第几周、星期几，构造一条必然命中的规则
	wom := _clocks.WeekOfMonthOf(now)
	weekday := now.Weekday()
	// 设为 1 秒后触发，避免刚好错过
	triggerTime := now.Add(2 * _time.Second)

	var count int32

	tasker := _clocks.TaskWeek.New(loc).
		AddRule(wom, weekday, triggerTime.Hour(), triggerTime.Minute(), triggerTime.Second()).
		SetName("周规则匹配测试").
		SetFn(func(_ _clocks.Tasker) {
			atomic.AddInt32(&count, 1)
			t.Logf("规则命中！当前时间: %s, 第%d周 %s", now.Format("2006-01-02 15:04:05"), wom, weekday)
		})

	errCh := make(chan error, 1)
	go func() { errCh <- tasker.Begin() }()

	// 等待足够时间让规则触发
	_time.Sleep(5 * _time.Second)
	tasker.Stop()

	if err := <-errCh; err != nil {
		t.Fatalf("任务执行失败：%v", err)
	}

	if atomic.LoadInt32(&count) != 1 {
		t.Fatalf("期望触发 1 次，实际触发 %d 次", count)
	}

	t.Logf("测试通过，触发 %d 次", count)
}

// Test_WeekRuleNoMatch 验证不匹配的规则不会触发：设置一个不可能命中的规则
func Test_WeekRuleNoMatch(t *_testing.T) {
	loc := _time.UTC

	var count int32

	// 构造一条当前时刻绝不会命中的规则：第5周星期一的 00:00:00（大多数月份第5周不存在）
	tasker := _clocks.TaskWeek.New(loc).
		AddRule(5, _time.Monday, 0, 0, 0).
		SetName("不匹配规则测试").
		SetFn(func(_ _clocks.Tasker) {
			atomic.AddInt32(&count, 1)
		})

	errCh := make(chan error, 1)
	go func() { errCh <- tasker.Begin() }()

	_time.Sleep(3 * _time.Second)
	tasker.Stop()

	if err := <-errCh; err != nil {
		t.Fatalf("任务执行失败：%v", err)
	}

	if atomic.LoadInt32(&count) != 0 {
		t.Fatalf("期望触发 0 次，实际触发 %d 次", count)
	}

	t.Log("测试通过，未命中规则没有触发")
}

// Test_WeekMultiRules 验证多条规则同时配置，只有匹配的规则会触发
func Test_WeekMultiRules(t *_testing.T) {
	loc, _ := _time.LoadLocation("Asia/Shanghai")
	now := _time.Now().In(loc)

	wom := _clocks.WeekOfMonthOf(now)
	weekday := now.Weekday()
	triggerTime := now.Add(2 * _time.Second)

	var count int32

	tasker := _clocks.TaskWeek.New(loc).
		// 命中的规则
		AddRule(wom, weekday, triggerTime.Hour(), triggerTime.Minute(), triggerTime.Second()).
		// 不命中的规则：第5周周五 03:00:00
		AddRule(5, _time.Friday, 3, 0, 0).
		SetName("多规则测试").
		SetFn(func(tasker _clocks.Tasker) {
			atomic.AddInt32(&count, 1)
			t.Logf("[%s] 规则命中，当前时间: %s", tasker.Name(), _time.Now().In(loc).Format("2006-01-02 15:04:05"))
		})

	errCh := make(chan error, 1)
	go func() { errCh <- tasker.Begin() }()

	_time.Sleep(5 * _time.Second)
	tasker.Stop()

	if err := <-errCh; err != nil {
		t.Fatalf("任务执行失败：%v", err)
	}

	if atomic.LoadInt32(&count) != 1 {
		t.Fatalf("期望触发 1 次，实际触发 %d 次", count)
	}

	t.Logf("测试通过，多规则场景触发 %d 次", count)
}

// Test_WeekNoRules 验证未配置规则时 Begin 返回错误
func Test_WeekNoRules(t *_testing.T) {
	tasker := _clocks.TaskWeek.New(nil).
		SetFn(func(_ _clocks.Tasker) {})

	err := tasker.Begin()
	if err == nil {
		t.Fatal("期望返回错误，但未返回")
	}
	t.Logf("符合预期： %v", err)
}

// Test_WeekNoFn 验证未设置回调时 Begin 返回错误
func Test_WeekNoFn(t *_testing.T) {
	tasker := _clocks.TaskWeek.New(nil).
		AddRule(1, _time.Monday, 12, 0, 0)

	err := tasker.Begin()
	if err == nil {
		t.Fatal("期望返回错误，但未返回")
	}
	t.Logf("符合预期：%v", err)
}
