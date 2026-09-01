package clocks_test

import (
	_testing "testing"
	_time "time"

	_clocks "github.com/aid297/aid/v2/clocks"
)

func newClock() _clocks.Clock { return _clocks.OnceClock() }

func Test_Clock_AfterNormal(t *_testing.T) {
	clock := newClock()
	defer clock.Clean()

	tasker := _clocks.NewTaskOnceAfter(3 * _time.Second).
		SetFn(func(_ _clocks.Tasker) { t.Logf("执行定时任务") }).
		SetName("定时任务").
		SetTimeout(5 * _time.Second)

	clock.
		AddTasker(tasker).
		SetErrHandler(func(tasker _clocks.Tasker, err error) { t.Errorf("定时任务：%v 执行失败：%v", tasker, err) }).
		Boot()

	_time.Sleep(4 * _time.Second)
	t.Log("测试完成")
}

func Test_Clock_AfterTimeout(t *_testing.T) {
	clock := newClock()
	defer clock.Clean()

	tasker1 := _clocks.NewTaskOnceAfter(5 * _time.Second).
		SetFn(func(_ _clocks.Tasker) { t.Logf("执行定时任务（会超时） -> 执行成功") }).
		SetName("定时任务1：会超时").
		SetTimeout(3 * _time.Second)

	tasker2 := _clocks.NewTaskOnceAfter(3 * _time.Second).
		SetFn(func(_ _clocks.Tasker) { t.Logf("执行定时任务（不会超时） ->  执行成功") }).
		SetName("定时任务2：不会超时").
		SetTimeout(5 * _time.Second)

	clock.
		AddTasker(tasker1, tasker2).
		SetErrHandler(func(tasker _clocks.Tasker, err error) { t.Errorf("定时任务：%v 执行失败：%v", tasker, err) }).
		Boot()

	_time.Sleep(4 * _time.Second)
	t.Log("测试完成")
}
