package clocks_test

import (
	_testing "testing"
	_time "time"

	_clocks "github.com/aid297/aid/v2/clocks"
)

func Test_CyclicityNormal(t *_testing.T) {
	count := 0

	tasker := _clocks.TaskCyclicity.
		Secondly(3).
		SetTimeout(5 * _time.Second).
		SetFn(func() { count++; t.Logf("每3秒执行一次：执行成功；第%d次", count) })

	errCh := make(chan error, 1)
	go func() {
		errCh <- tasker.Begin()
	}()

	// 等待 10 秒后停止任务
	_time.Sleep(10 * _time.Second)
	tasker.Stop()

	if err := <-errCh; err != nil {
		t.Fatalf("任务（%v），执行失败：%v", tasker, err)
	}

	t.Logf("测试结束，共执行 %d 次：%v", count, tasker)
}

func Test_CyclicityTimeout(t *_testing.T) {
	clock := _clocks.Clock.Ins().SetErrHandler(func(tasker _clocks.Tasker, err error) { t.Errorf("执行错误(%v)：%v", tasker, err) })
	defer clock.Clean()

	count := 0

	tasker1 := _clocks.TaskCyclicity.
		Secondly(3).
		SetName("任务A").
		SetFn(func() { count++; t.Logf("【任务A】每3秒执行一次，执行成功：第%d次", count) })

	tasker2 := _clocks.TaskCyclicity.Secondly(3).SetName("任务B").SetFn(func() { count++; t.Logf("【任务B】每3秒执行一次，执行成功：第%d次", count) })

	clock.AddTasker(tasker1)
	clock.Boot()

	clock.AddTaskerAndBegin(tasker2)

	// 等待 10 秒后停止任务
	_time.Sleep(10 * _time.Second)

	t.Logf("测试结束，共执行 %d 次：%v", count, tasker1)
}
