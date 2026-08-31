package retries_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aid297/aid/v2/debugLogs"
	"github.com/aid297/aid/v2/retries"
)

func operations() error {
	debugLogs.Print("Executing operations...")
	return errors.New("transient error")
}

func Test1(t *testing.T) {
	t.Run("test1 指数退避重试", func(t *testing.T) {
		err := retries.APP.Retry.New(retries.Sleep(time.Second), retries.Fn(operations)).Exponent(3)
		if err != nil {
			t.Logf("Operation failed after retries: %v", err)
		}
	})
}

func Test2(t *testing.T) {
	t.Run("test2 支持上下文的匀速重试", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err := retries.APP.Retry.New(retries.Sleep(time.Second), retries.Fn(operations), retries.Context(ctx)).LinearWithContext(3)
		if err != nil {
			t.Logf("Operation failed after retries: %v", err)
		}
	})
}

func Test3(t *testing.T) {
	t.Run("test3 支持上下文的随机退避重试", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err := retries.APP.Retry.New(retries.Sleep(time.Second), retries.Fn(operations), retries.Context(ctx)).JitterWithContext(5)
		if err != nil {
			t.Logf("Operation failed after retries: %v", err)
		}
	})
}
