package clocks_test

import (
	_syncAtomic "sync/atomic"
	_testing "testing"
	_time "time"

	_clocks "github.com/aid297/aid/v2/clocks"
)

func Test_CyclicityNormal(t *_testing.T) {
	count := 0

	tasker := _clocks.NewTaskCyclicitySecondly(3).
		SetTimeout(5 * _time.Second).
		SetHandler(func(_ _clocks.Tasker) { count++; t.Logf("每3秒执行一次：执行成功；第%d次", count) })

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
	clock := _clocks.OnceClock().SetErrHandler(func(tasker _clocks.Tasker, err error) { t.Errorf("执行错误(%v)：%v", tasker, err) })
	defer clock.Clean()

	var count int32

	tasker1 := _clocks.NewTaskCyclicitySecondly(3).
		SetName("任务A").
		SetHandler(func(_ _clocks.Tasker) {
			_syncAtomic.AddInt32(&count, 1)
			t.Logf("【任务A】每3秒执行一次，执行成功：第%d次", _syncAtomic.LoadInt32(&count))
		})

	tasker2 := _clocks.NewTaskCyclicitySecondly(3).SetName("任务B").SetHandler(func(_ _clocks.Tasker) {
		_syncAtomic.AddInt32(&count, 1)
		t.Logf("【任务B】每3秒执行一次，执行成功：第%d次", _syncAtomic.LoadInt32(&count))
	})

	clock.AddTasker(tasker1)
	clock.Boot()

	clock.AddTaskerAndBegin(tasker2)

	// 等待 10 秒后停止任务
	_time.Sleep(10 * _time.Second)

	t.Logf("测试结束，共执行 %d 次：%v", _syncAtomic.LoadInt32(&count), tasker1)
}
