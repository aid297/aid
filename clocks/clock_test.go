package clocks

import (
	_testing "testing"
	_time "time"
)

func newClock() ClockImpl { return Clock.Ins() }

func Test_Clock_AfterNormal(t *_testing.T) {
	clock := newClock()
	defer clock.Clean()

	tasker := TaskOnce.
		After(3 * _time.Second).
		SetFn(func() { t.Logf("执行定时任务") }).
		SetName("定时任务").
		SetTimeout(5 * _time.Second)

	clock = clock.
		AddTasker(tasker).
		SetErrHandler(func(tasker Tasker, err error) { t.Errorf("定时任务：%v 执行失败：%v", tasker, err) }).
		Boot()

	_time.Sleep(4 * _time.Second)
	t.Log("测试完成")
}

func Test_Clock_AfterTimeout(t *_testing.T) {
	clock := newClock()
	defer clock.Clean()

	tasker1 := TaskOnce.
		After(5 * _time.Second).
		SetFn(func() { t.Logf("执行定时任务（会超时） -> 执行成功") }).
		SetName("定时任务1：会超时").
		SetTimeout(3 * _time.Second)

	tasker2 := TaskOnce.
		After(3 * _time.Second).
		SetFn(func() { t.Logf("执行定时任务（不会超时） ->  执行成功") }).
		SetName("定时任务2：不会超时").
		SetTimeout(5 * _time.Second)

	clock.
		AddTasker(tasker1, tasker2).
		SetErrHandler(func(tasker Tasker, err error) { t.Errorf("定时任务：%v 执行失败：%v", tasker, err) }).
		Boot()

	_time.Sleep(4 * _time.Second)
	t.Log("测试完成")
}
