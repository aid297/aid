package clocks_test

import (
	"sync/atomic"
	_testing "testing"
	_time "time"

	_clocks "github.com/aid297/aid/v2/clocks"
)

// Test_WeekOfMonthOf 验证 weekOfMonthOf 工具函数的周序号计算
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

	tasker := _clocks.NewTaskWeek(loc).
		AddRule(_clocks.EveryMonth, wom, weekday, triggerTime.Hour(), triggerTime.Minute(), triggerTime.Second()).
		SetName("周规则匹配测试").
		SetHandler(func(_ _clocks.Tasker) {
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
	tasker := _clocks.NewTaskWeek(loc).
		AddRule(_clocks.EveryMonth, 5, _time.Monday, 0, 0, 0).
		SetName("不匹配规则测试").
		SetHandler(func(_ _clocks.Tasker) {
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

	tasker := _clocks.NewTaskWeek(loc).
		// 命中的规则
		AddRule(_clocks.EveryMonth, wom, weekday, triggerTime.Hour(), triggerTime.Minute(), triggerTime.Second()).
		// 不命中的规则：第5周周五 03:00:00
		AddRule(_clocks.EveryMonth, 5, _time.Friday, 3, 0, 0).
		SetName("多规则测试").
		SetHandler(func(tasker _clocks.Tasker) {
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

// Test_LastWeekOfMonth 验证 lastWeekOfMonth 工具函数：不同月份天数不同，最后一周可能是4或5
func Test_LastWeekOfMonth(t *_testing.T) {
	loc := _time.UTC
	cases := []struct {
		date _time.Time
		want int
		desc string
	}{
		{_time.Date(2026, 2, 15, 0, 0, 0, 0, loc), 4, "2月（28天）最后一周为第4周"},
		{_time.Date(2026, 8, 15, 0, 0, 0, 0, loc), 5, "8月（31天）最后一周为第5周"},
		{_time.Date(2026, 4, 15, 0, 0, 0, 0, loc), 5, "4月（30天）最后一周为第5周"},
		{_time.Date(2026, 11, 15, 0, 0, 0, 0, loc), 5, "11月（30天）最后一周为第5周"},
	}

	for _, c := range cases {
		got := _clocks.LastWeekOfMonth(c.date)
		if got != c.want {
			t.Errorf("LastWeekOfMonth(%s) = %d, want %d (%s)", c.date.Format("2006-01"), got, c.want, c.desc)
		} else {
			t.Logf("✓ %s: LastWeekOfMonth = %d", c.desc, got)
		}
	}
}

// Test_WeekLastWeekRule 验证 LastWeek 规则在当前月份最后一周能正确命中
func Test_WeekLastWeekRule(t *_testing.T) {
	loc, _ := _time.LoadLocation("Asia/Shanghai")
	now := _time.Now().In(loc)
	lastWom := _clocks.LastWeekOfMonth(now)
	currentWom := _clocks.WeekOfMonthOf(now)

	// 仅在当前周就是最后一周时运行此测试
	if currentWom != lastWom {
		t.Skipf("当前是第%d周，本月最后一周是第%d周，跳过实时匹配测试", currentWom, lastWom)
		return
	}

	triggerTime := now.Add(2 * _time.Second)
	var count int32

	tasker := _clocks.NewTaskWeek(loc).
		AddRule(_clocks.EveryMonth, _clocks.LastWeek, now.Weekday(), triggerTime.Hour(), triggerTime.Minute(), triggerTime.Second()).
		SetName("LastWeek 规则测试").
		SetHandler(func(_ _clocks.Tasker) {
			atomic.AddInt32(&count, 1)
			t.Logf("LastWeek 规则命中！当前月份最后一周 = %d, 当前时间: %s", lastWom, now.Format("2006-01-02 15:04:05"))
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

	t.Logf("测试通过，LastWeek 正确解析为第%d周", lastWom)
}

// Test_WeekNoRules 验证未配置规则时 Begin 返回错误
func Test_WeekNoRules(t *_testing.T) {
	tasker := _clocks.NewTaskWeek(nil).
		SetHandler(func(_ _clocks.Tasker) {})

	err := tasker.Begin()
	if err == nil {
		t.Fatal("期望返回错误，但未返回")
	}
	t.Logf("符合预期： %v", err)
}

// Test_WeekNoFn 验证未设置回调时 Begin 返回错误
func Test_WeekNoFn(t *_testing.T) {
	tasker := _clocks.NewTaskWeek(nil).
		AddRule(_clocks.EveryMonth, 1, _time.Monday, 12, 0, 0)

	err := tasker.Begin()
	if err == nil {
		t.Fatal("期望返回错误，但未返回")
	}
	t.Logf("符合预期：%v", err)
}

// Test_WeekMonthFilter 验证月份过滤：指定非当前月份的规则不会触发
func Test_WeekMonthFilter(t *_testing.T) {
	loc, _ := _time.LoadLocation("Asia/Shanghai")
	now := _time.Now().In(loc)

	// 构造一个当前月不匹配的月份
	wrongMonth := now.Month() + 1
	if wrongMonth > _time.December {
		wrongMonth = now.Month() - 1
	}

	wom := _clocks.WeekOfMonthOf(now)
	triggerTime := now.Add(2 * _time.Second)

	var count int32

	tasker := _clocks.NewTaskWeek(loc).
		// 指定为其他月份，即使周、星期、时刻都匹配也不应触发
		AddRule(wrongMonth, wom, now.Weekday(), triggerTime.Hour(), triggerTime.Minute(), triggerTime.Second()).
		SetName("月份过滤测试").
		SetHandler(func(_ _clocks.Tasker) {
			atomic.AddInt32(&count, 1)
		})

	errCh := make(chan error, 1)
	go func() { errCh <- tasker.Begin() }()

	_time.Sleep(5 * _time.Second)
	tasker.Stop()

	if err := <-errCh; err != nil {
		t.Fatalf("任务执行失败：%v", err)
	}

	if atomic.LoadInt32(&count) != 0 {
		t.Fatalf("期望触发 0 次（当前 %d 月，规则设置为 %d 月），实际触发 %d 次", now.Month(), wrongMonth, count)
	}

	t.Logf("测试通过，月份过滤生效（当前 %d 月，规则 %d 月未触发）", now.Month(), wrongMonth)
}

// Test_WeekSpecificMonthMatch 验证指定当前月份的规则能正确触发
func Test_WeekSpecificMonthMatch(t *_testing.T) {
	loc, _ := _time.LoadLocation("Asia/Shanghai")
	now := _time.Now().In(loc)

	wom := _clocks.WeekOfMonthOf(now)
	triggerTime := now.Add(2 * _time.Second)
	var count int32

	tasker := _clocks.NewTaskWeek(loc).
		// 指定当前月份
		AddRule(now.Month(), wom, now.Weekday(), triggerTime.Hour(), triggerTime.Minute(), triggerTime.Second()).
		SetName("指定月份匹配测试").
		SetHandler(func(_ _clocks.Tasker) {
			atomic.AddInt32(&count, 1)
			t.Logf("指定月份规则命中！当前 %d 月，第%d周 %s", now.Month(), wom, now.Weekday())
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

	t.Logf("测试通过，指定当前月份规则触发 %d 次", count)
}
