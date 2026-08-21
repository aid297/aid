### RedisPool 说明

```go
package redisPool_test

import (
	"context"
	"testing"
	"time"

	"github.com/aid297/aid/v3/redisPools"
)

func TestOnceRedisPool(t *testing.T) {
	rp := redisPools.New(
		"localhost:6379",
		"business",
		redisPools.Pool( "users", 0),
		redisPools.Pool( "orders", 1),
	)

	prefix, rc := rp.GetClient("user")
	if rc == nil {
		t.Fatal("user client is nil")
	}

	if prefix != "business:users" {
		t.Fatal("user prefix is not correct")
	}

	if res := rc.Set(context.Background(), "some-key", "some-value", 5*time.Minute); res.Err() != nil {
		t.Fatal("set user key error: ", res.Err())
	}

	if res := rc.Get(context.Background(), "some-key"); res.Err() != nil {
		t.Fatal("get user key error: ", res.Err())
	}

	rp.Set(context.Background(), "user", "some-key2", "some-value2", 5*time.Minute)
	if someValue2, err := rp.Get(context.Background(), "user", "some-key2"); err != nil {
		t.Fatal("get user key2 error: ", err)
	} else {
		t.Log("some-value2: ", someValue2)
	}

	t.Cleanup(func() { rp.Clean() })
}

```

