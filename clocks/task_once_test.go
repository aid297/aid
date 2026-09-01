package clocks_test

import (
	_testing "testing"
	_time "time"

	_clocks "github.com/aid297/aid/v2/clocks"
)

func Test_AfterNormal(t *_testing.T) {
	if err := _clocks.NewTaskOnceAfter(3 * _time.Second).
		SetTimeout(5 * _time.Second).
		SetFn(func(_ _clocks.Tasker) { t.Logf("执行定时任务") }).
		Begin(); err != nil {
		t.Fatalf("测试失败：%v", err)
	}

	t.Logf("测试完成")
}

func Test_AfterTimeout(t *_testing.T) {
	err := _clocks.NewTaskOnceAfter(3 * _time.Second).
		SetTimeout(1 * _time.Second).
		SetFn(func(_ _clocks.Tasker) { t.Logf("执行定时任务") }).
		Begin()
	if err != nil {
		t.Errorf("测试失败：%v", err)
	}

	t.Logf("测试结束")
}

func Test_AtNormal(t *_testing.T) {
	loc, err := _time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatalf("设置时区错误：%v", err)
	}
	if err = _clocks.NewTaskOnceAt(_time.Now().In(loc).Add(3*_time.Second), loc).
		SetTimeout(5 * _time.Second).
		SetFn(func(_ _clocks.Tasker) { t.Logf("执行定时任务") }).
		Begin(); err != nil {
		t.Errorf(" 测试失败：%v", err)
	}

	t.Logf("测试结束")
}

func Test_AtTimeout(t *_testing.T) {
	loc, err := _time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatalf("设置时区错误：%v", err)
	}
	if err = _clocks.NewTaskOnceAt(_time.Now().In(loc).Add(3*_time.Second), loc).
		SetTimeout(1 * _time.Second).
		SetFn(func(_ _clocks.Tasker) { t.Logf("执行定时任务") }).
		Begin(); err != nil {
		t.Errorf(" 测试失败：%v", err)
	}

	t.Logf("测试结束")
}
