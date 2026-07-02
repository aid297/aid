package coroutineGroupV2

import (
	"math/rand"
	"time"
)

// BackoffStrategy 退避策略
type BackoffStrategy uint8

const (
	// BackoffFixed 固定间隔退避：每次重试等待 interval
	BackoffFixed BackoffStrategy = iota
	// BackoffExponential 指数退避：第n次重试等待 interval * 2^(n-1)
	BackoffExponential
	// BackoffExponentialJitter 指数退避+随机抖动：在指数退避基础上加入随机抖动
	BackoffExponentialJitter
)

// RetryConfig 重试配置
type RetryConfig struct {
	maxAttempts int           // 最大执行次数（含首次，例如3=首次+2次重试）
	timeout     time.Duration // 每次执行的超时时间，0表示不超时
	interval    time.Duration // 重试基础间隔时间
	backoff     BackoffStrategy
}

// CoroutineGroupRetryAttr 重试配置选项
type CoroutineGroupRetryAttr func(*RetryConfig)

// WithRetryAttempts 设置最大执行次数（含首次，例如3=首次+2次重试）
func WithRetryAttempts(attempts int) CoroutineGroupRetryAttr {
	return func(c *RetryConfig) {
		if attempts < 1 {
			attempts = 1
		}
		c.maxAttempts = attempts
	}
}

// WithTimeout 设置每次执行的超时时间，0表示不超时
func WithTimeout(timeout time.Duration) CoroutineGroupRetryAttr {
	return func(c *RetryConfig) {
		c.timeout = timeout
	}
}

// WithRetryInterval 设置重试基础间隔时间
func WithRetryInterval(interval time.Duration) CoroutineGroupRetryAttr {
	return func(c *RetryConfig) {
		c.interval = interval
	}
}

// WithBackoff 设置退避策略
func WithBackoff(strategy BackoffStrategy) CoroutineGroupRetryAttr {
	return func(c *RetryConfig) {
		c.backoff = strategy
	}
}

// calculateWait 计算第 attempt 次重试的等待时间（attempt 从1开始）
func (c *RetryConfig) calculateWait(attempt int) time.Duration {
	if c.interval <= 0 {
		return 0
	}

	switch c.backoff {
	case BackoffExponential:
		// interval * 2^(attempt-1)
		return c.interval * time.Duration(1<<(attempt-1))
	case BackoffExponentialJitter:
		// interval * 2^(attempt-1) + random[0, base]
		base := c.interval * time.Duration(1<<(attempt-1))
		jitter := time.Duration(rand.Int63n(int64(base) + 1))
		return base + jitter
	default: // BackoffFixed
		return c.interval
	}
}
