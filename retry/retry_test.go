package retry

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func operation() error {
	fmt.Println("Executing operation...")
	return errors.New("transient error")
}

func Test1(t *testing.T) {
	t.Run("test1 指数退避重试", func(t *testing.T) {
		err := APP.Retry.New(Sleep(time.Second), Fn(operation)).Do(3)
		if err != nil {
			t.Logf("Operation failed after retries: %v", err)
		}
	})
}

func Test2(t *testing.T) {
	t.Run("test2 支持上下文的匀速重试", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err := APP.Retry.New(Sleep(time.Second), Fn(operation), Context(ctx)).WithContext(3)
		if err != nil {
			t.Logf("Operation failed after retries: %v", err)
		}
	})
}

func Test3(t *testing.T) {
	t.Run("test3 支持上下文的随机退避重试", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err := APP.Retry.New(Sleep(time.Second), Fn(operation), Context(ctx)).WithContextAndJitter(5)
		if err != nil {
			t.Logf("Operation failed after retries: %v", err)
		}
	})
}
