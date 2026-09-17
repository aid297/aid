package retries_test

import (
	_context "context"
	_errors "errors"
	_testing "testing"
	_time "time"

	_retries "github.com/aid297/aid/v2/retries"
)

// errMock 用于测试重试流程中函数返回的错误
var errMock = _errors.New("mock error")

// Test_NewRetry 验证构造函数返回可用的 Retry 实例
func Test_NewRetry(t *_testing.T) {
	if r := _retries.NewRetry(); r == nil {
		t.Fatal("NewRetry 应返回非 nil 的 Retry 实例")
	}
}

// Test_Set_Chainable 验证 Set 返回自身以支持链式调用
func Test_Set_Chainable(t *_testing.T) {
	called := 0
	r := _retries.NewRetry().
		Set(_retries.Sleep(1 * _time.Millisecond)).
		Set(_retries.Fn(func(retried uint) error {
			called++
			return nil
		}))

	if err := r.Linear(1); err != nil {
		t.Fatalf("期望成功，实际错误：%v", err)
	}
	if called != 1 {
		t.Fatalf("期望调用 1 次，实际 %d 次", called)
	}
}

// Test_Linear_NilFn 未设置 fn 时应直接返回 nil
func Test_Linear_NilFn(t *_testing.T) {
	if err := _retries.NewRetry().Linear(3); err != nil {
		t.Fatalf("期望返回 nil，实际：%v", err)
	}
}

// Test_Linear_FirstAttemptSuccess 首次调用即成功
func Test_Linear_FirstAttemptSuccess(t *_testing.T) {
	called := 0
	r := _retries.NewRetry(
		_retries.Sleep(1*_time.Millisecond),
		_retries.Fn(func(retried uint) error {
			called++
			return nil
		}),
	)

	if err := r.Linear(3); err != nil {
		t.Fatalf("期望成功，实际错误：%v", err)
	}
	if called != 1 {
		t.Fatalf("期望调用 1 次，实际 %d 次", called)
	}
}

// Test_Linear_SuccessAfterRetries 若干次失败后成功
func Test_Linear_SuccessAfterRetries(t *_testing.T) {
	called := 0
	r := _retries.NewRetry(
		_retries.Sleep(1*_time.Millisecond),
		_retries.Fn(func(retried uint) error {
			called++
			if called < 3 {
				return errMock
			}
			return nil
		}),
	)

	if err := r.Linear(5); err != nil {
		t.Fatalf("期望成功，实际错误：%v", err)
	}
	if called != 3 {
		t.Fatalf("期望调用 3 次，实际 %d 次", called)
	}
}

// Test_Linear_AllAttemptsFail 全部尝试均失败时应返回最后一次错误
func Test_Linear_AllAttemptsFail(t *_testing.T) {
	called := 0
	r := _retries.NewRetry(
		_retries.Sleep(1*_time.Millisecond),
		_retries.Fn(func(retried uint) error {
			called++
			return errMock
		}),
	)

	err := r.Linear(3)
	if !_errors.Is(err, errMock) {
		t.Fatalf("期望返回 errMock，实际：%v", err)
	}
	if called != 3 {
		t.Fatalf("期望调用 3 次，实际 %d 次", called)
	}
}

// Test_Exponent_NilFn 未设置 fn 时应直接返回 nil
func Test_Exponent_NilFn(t *_testing.T) {
	if err := _retries.NewRetry().Exponent(3); err != nil {
		t.Fatalf("期望返回 nil，实际：%v", err)
	}
}

// Test_Exponent_FirstAttemptSuccess 首次调用即成功
func Test_Exponent_FirstAttemptSuccess(t *_testing.T) {
	called := 0
	r := _retries.NewRetry(
		_retries.Sleep(1*_time.Millisecond),
		_retries.Fn(func(retried uint) error {
			called++
			return nil
		}),
	)

	if err := r.Exponent(3); err != nil {
		t.Fatalf("期望成功，实际错误：%v", err)
	}
	if called != 1 {
		t.Fatalf("期望调用 1 次，实际 %d 次", called)
	}
}

// Test_Exponent_SuccessAfterRetries 指数退避过程中成功
func Test_Exponent_SuccessAfterRetries(t *_testing.T) {
	called := 0
	r := _retries.NewRetry(
		_retries.Sleep(1*_time.Millisecond),
		_retries.Fn(func(retried uint) error {
			called++
			if called < 2 {
				return errMock
			}
			return nil
		}),
	)

	if err := r.Exponent(3); err != nil {
		t.Fatalf("期望成功，实际错误：%v", err)
	}
	if called != 2 {
		t.Fatalf("期望调用 2 次，实际 %d 次", called)
	}
}

// Test_Exponent_AllAttemptsFail 指数退避全部失败
func Test_Exponent_AllAttemptsFail(t *_testing.T) {
	called := 0
	r := _retries.NewRetry(
		_retries.Sleep(1*_time.Millisecond),
		_retries.Fn(func(retried uint) error {
			called++
			return errMock
		}),
	)

	err := r.Exponent(3)
	if !_errors.Is(err, errMock) {
		t.Fatalf("期望返回 errMock，实际：%v", err)
	}
	if called != 3 {
		t.Fatalf("期望调用 3 次，实际 %d 次", called)
	}
}

// Test_LinearWithContext_NilFn 未设置 fn 时应直接返回 nil
func Test_LinearWithContext_NilFn(t *_testing.T) {
	if err := _retries.NewRetry().LinearWithContext(3); err != nil {
		t.Fatalf("期望返回 nil，实际：%v", err)
	}
}

// Test_LinearWithContext_Success 上下文未取消时正常成功
func Test_LinearWithContext_Success(t *_testing.T) {
	called := 0
	r := _retries.NewRetry(
		_retries.Sleep(1*_time.Millisecond),
		_retries.Context(_context.Background()),
		_retries.Fn(func(retried uint) error {
			called++
			if called < 2 {
				return errMock
			}
			return nil
		}),
	)

	if err := r.LinearWithContext(3); err != nil {
		t.Fatalf("期望成功，实际错误：%v", err)
	}
	if called != 2 {
		t.Fatalf("期望调用 2 次，实际 %d 次", called)
	}
}

// Test_LinearWithContext_AllAttemptsFail 上下文未取消但全部失败
func Test_LinearWithContext_AllAttemptsFail(t *_testing.T) {
	r := _retries.NewRetry(
		_retries.Sleep(1*_time.Millisecond),
		_retries.Context(_context.Background()),
		_retries.Fn(func(retried uint) error {
			return errMock
		}),
	)

	err := r.LinearWithContext(2)
	if !_errors.Is(err, errMock) {
		t.Fatalf("期望返回 errMock，实际：%v", err)
	}
}

// Test_LinearWithContext_ContextCanceled 上下文被取消时应返回 ctx.Err()
func Test_LinearWithContext_ContextCanceled(t *_testing.T) {
	ctx, cancel := _context.WithCancel(_context.Background())
	r := _retries.NewRetry(
		_retries.Sleep(50*_time.Millisecond),
		_retries.Context(ctx),
		_retries.Fn(func(retried uint) error {
			// 首次调用后立即取消上下文，进入等待分支时应立刻感知取消
			cancel()
			return errMock
		}),
	)

	err := r.LinearWithContext(3)
	if !_errors.Is(err, _context.Canceled) {
		t.Fatalf("期望返回 context.Canceled，实际：%v", err)
	}
}

// Test_JitterWithContext_NilFn 未设置 fn 时应直接返回 nil
func Test_JitterWithContext_NilFn(t *_testing.T) {
	if err := _retries.NewRetry().JitterWithContext(3); err != nil {
		t.Fatalf("期望返回 nil，实际：%v", err)
	}
}

// Test_JitterWithContext_Success 抖动退避过程中成功
func Test_JitterWithContext_Success(t *_testing.T) {
	called := 0
	r := _retries.NewRetry(
		_retries.Sleep(2*_time.Millisecond),
		_retries.Context(_context.Background()),
		_retries.Fn(func(retried uint) error {
			called++
			if called < 2 {
				return errMock
			}
			return nil
		}),
	)

	if err := r.JitterWithContext(3); err != nil {
		t.Fatalf("期望成功，实际错误：%v", err)
	}
	if called != 2 {
		t.Fatalf("期望调用 2 次，实际 %d 次", called)
	}
}

// Test_JitterWithContext_AllAttemptsFail 抖动退避全部失败
func Test_JitterWithContext_AllAttemptsFail(t *_testing.T) {
	r := _retries.NewRetry(
		_retries.Sleep(1*_time.Millisecond),
		_retries.Context(_context.Background()),
		_retries.Fn(func(retried uint) error {
			return errMock
		}),
	)

	err := r.JitterWithContext(2)
	if !_errors.Is(err, errMock) {
		t.Fatalf("期望返回 errMock，实际：%v", err)
	}
}

// Test_JitterWithContext_ContextCanceled 抖动退避中上下文被取消
func Test_JitterWithContext_ContextCanceled(t *_testing.T) {
	ctx, cancel := _context.WithCancel(_context.Background())
	r := _retries.NewRetry(
		_retries.Sleep(50*_time.Millisecond),
		_retries.Context(ctx),
		_retries.Fn(func(retried uint) error {
			cancel()
			return errMock
		}),
	)

	err := r.JitterWithContext(3)
	if !_errors.Is(err, _context.Canceled) {
		t.Fatalf("期望返回 context.Canceled，实际：%v", err)
	}
}
