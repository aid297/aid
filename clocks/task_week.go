package clocks

import (
	_context "context"
	_errors "errors"
	"fmt"
	_time "time"

	_uuid "github.com/google/uuid"
)

var (
	_        Tasker = (*TaskWeekImpl)(nil)
	TaskWeek TaskWeekImpl
)

const (
	// LastWeek 表示月份最后一周：在只有4周的月份等价于第4周，在有5周的月份等价于第5周
	LastWeek = -1
)

var (
	// EveryMonth 表示每月都生效（不限制月份），用于 AddRule 的 month 参数
	EveryMonth = _time.Month(0)
)

// WeekRule 定义月份中第 N 周的某个星期几在指定时刻触发的规则
//
// Month: 指定生效月份（_time.January ~ _time.December）；EveryMonth（0）表示每月都生效
// WeekOfMonth: 1-5 表示月份的第几周（第1周=1-7日，第2周=8-14日，以此类推）；
// 传 LastWeek（-1）表示月份最后一周，自动适配当月实际周数（4或5）
// Weekday: 星期几（_time.Sunday ~ _time.Saturday）
// Hour/Minute/Second: 触发时刻
type WeekRule struct {
	Month       _time.Month   // 生效月份（1-12），EveryMonth（0）表示每月
	WeekOfMonth int           // 第几周（1-5），LastWeek（-1）表示最后一周
	Weekday     _time.Weekday // 星期几
	Hour        int           // 时（0-23）
	Minute      int           // 分（0-59）
	Second      int           // 秒（0-59）
}

// WeekOfMonthOf 根据日期计算当前是月份的第几周（第1-7日为第1周，第8-14日为第2周，以此类推）
func WeekOfMonthOf(t _time.Time) int {
	return (t.Day()-1)/7 + 1
}

// LastWeekOfMonth 计算指定日期所在月份共有几周（返回 4 或 5）
func LastWeekOfMonth(t _time.Time) int {
	lastDay := _time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location())
	return WeekOfMonthOf(lastDay)
}

type TaskWeekImpl struct {
	uuid        _uuid.UUID
	name        string
	timeout     _time.Duration
	fn          func(tasker Tasker)
	closeCh     chan es
	loc         *_time.Location
	rules       []WeekRule
	immediately bool
	lastFired   _time.Time // 上次触发的时刻（精确到秒），用于防止同一秒内重复触发
}

// New 创建按周规则定时任务，loc 为时区（传 nil 则使用 time.Local）
func (*TaskWeekImpl) New(loc *_time.Location) *TaskWeekImpl {
	if loc == nil {
		loc = _time.Local
	}
	return &TaskWeekImpl{
		uuid:    _uuid.Must(_uuid.NewV7()),
		timeout: defaultTimeout,
		closeCh: make(chan es, 1),
		loc:     loc,
	}
}

// AddRule 添加一条触发规则：month 月的第 weekOfMonth 周的 weekday 在 hour:minute:second 触发
//
// month 取值范围：_time.January ~ _time.December 或 EveryMonth（0）表示每月都生效
// weekOfMonth 取值范围：1-5 或 LastWeek（-1）表示最后一周
func (my *TaskWeekImpl) AddRule(month _time.Month, weekOfMonth int, weekday _time.Weekday, hour, minute, second int) *TaskWeekImpl {
	if month < _time.January || month > _time.December {
		month = EveryMonth
	}
	if weekOfMonth != LastWeek {
		if weekOfMonth < 1 {
			weekOfMonth = 1
		}
		if weekOfMonth > 5 {
			weekOfMonth = 5
		}
	}
	my.rules = append(my.rules, WeekRule{
		Month:       month,
		WeekOfMonth: weekOfMonth,
		Weekday:     weekday,
		Hour:        hour,
		Minute:      minute,
		Second:      second,
	})
	return my
}

func (my *TaskWeekImpl) String() string {
	return fmt.Sprintf("uuid: %s, name: %s, rules: %d, timeout: %s", my.uuid, my.name, len(my.rules), my.timeout)
}

func (my *TaskWeekImpl) UUID() _uuid.UUID { return my.uuid }

func (my *TaskWeekImpl) SetName(name string) Tasker { my.name = name; return my }

func (my *TaskWeekImpl) Name() string { return my.name }

func (my *TaskWeekImpl) SetTimeout(timeout _time.Duration) Tasker {
	if timeout > 0 {
		my.timeout = timeout
	}
	return my
}

func (my *TaskWeekImpl) Timeout() _time.Duration { return my.timeout }

func (my *TaskWeekImpl) SetFn(fn func(tasker Tasker)) Tasker { my.fn = fn; return my }

func (my *TaskWeekImpl) Fn() func(tasker Tasker) { return my.fn }

func (my *TaskWeekImpl) SetImmediately(immediately bool) Tasker {
	my.immediately = immediately
	return my
}

func (my *TaskWeekImpl) Immediately() bool { return my.immediately }

func (my *TaskWeekImpl) Do() {
	ctx, cancel := _context.WithTimeout(_context.Background(), my.timeout)
	defer cancel()

	done := make(chan struct{})
	go func() { defer close(done); my.fn(my) }()

	select {
	case <-ctx.Done():
		clockIns.errHandler(my, _errors.New("任务执行超时"))
	case <-done:
	}
}

// matchRule 检查当前时间是否命中某条规则，返回命中的规则指针（未命中返回 nil）
func (my *TaskWeekImpl) matchRule(now _time.Time) *WeekRule {
	wom := WeekOfMonthOf(now)
	lastWom := LastWeekOfMonth(now)
	curMonth := now.Month()
	for idx := range my.rules {
		r := &my.rules[idx]
		// 月份匹配：EveryMonth 表示不限制月份
		if r.Month != EveryMonth && r.Month != curMonth {
			continue
		}
		// LastWeek 动态解析为当月实际最后一周
		ruleWom := r.WeekOfMonth
		if ruleWom == LastWeek {
			ruleWom = lastWom
		}
		if ruleWom == wom &&
			r.Weekday == now.Weekday() &&
			r.Hour == now.Hour() &&
			r.Minute == now.Minute() &&
			r.Second == now.Second() {
			return r
		}
	}
	return nil
}

func (my *TaskWeekImpl) Begin() error {
	if my.fn == nil {
		return _errors.New("回调方法为空")
	}
	if len(my.rules) == 0 {
		return _errors.New("未配置任何周规则")
	}

	if my.Immediately() {
		my.Do()
	}

	ticker := _time.NewTicker(1 * _time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-my.closeCh:
			return nil
		case <-ticker.C:
			now := _time.Now().In(my.loc)
			// 截断到秒级，防止同一秒内多次触发
			nowTrunc := now.Truncate(1 * _time.Second)
			if !my.lastFired.IsZero() && !nowTrunc.After(my.lastFired) {
				continue
			}
			if r := my.matchRule(now); r != nil {
				my.lastFired = nowTrunc
				my.Do()
			}
		}
	}
}

func (my *TaskWeekImpl) Stop() Tasker { my.closeCh <- es{}; return my }
