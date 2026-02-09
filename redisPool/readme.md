### RedisPool 说明

```go
package main

import (
	`context`
	`fmt`
	`log`

	`github.com/aid297/aid/redisPool`
	`github.com/aid297/aid/time`
)

func main() {
	redisSetting, err := redisPool.Elements.Setting.New("redis.yaml")
	if err != nil {
		panic(err)
	}
	rp := redisPool.Elements.Pool.Once(redisSetting)
	prefix, rc := rp.GetClient("default")

	res := rc.Set(context.Background(), fmt.Sprintf("%s:%s", prefix, "some-key"), "some-value", 5*time.Minute)
	if err = res.Err(); err != nil {
		panic(err)
	}

	res2 := rc.Get(context.Background(), fmt.Sprintf("%s:%s", prefix, "some-key"))
	if err = res.Err(); err != nil {
		panic(err)
	}
	fmt.Printf("some-value: %s\n", res2.String())

	if _, err = rp.Set("default", "some-key2", "some-value2", 5*time.Minute); err != nil {
		panic(err)
	}

	someValue2, err := rp.Get("default", "some-key2")
	if err != nil {
		panic(err)
	}

	log.Printf("some-value2: %v", someValue2)
}
```

