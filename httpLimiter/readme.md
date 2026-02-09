### HTTPLimiter 限流器说明

1. `IP限流器`

   ```go
   package main
   
   import (
   	`fmt`
   	`time`
   
   	`github.com/aid297/aid/httpLimiter`
   )
   
   func main() {
   	v, ok := httpLimiter.APP.IPLimiter.New().Affirm("ip", 30*time.Second, 50)
   	if !ok {
   		panic(fmt.Errorf("[ip方式] %s 内访问了 %d 次", time.Since(v.GetLastVisitor()), v.GetVisitTimes()))
   	}
   }
   ```

2. `路由限流器`

   ```go
   package main
   
   import (
   	`fmt`
   	`time`
   
   	`github.com/aid297/aid/httpLimiter`
   )
   
   func main() {
   	rl := httpLimiter.APP.RouterLimiter.Once().
   		Add("/route1", 30*time.Second, 50).
   		Add("/route2", 30*time.Second, 100).
   		Add("/route3", 30*time.Second, 150)
   
   	v, ok := rl.Affirm("/route1", "ip")
   	if !ok {
   		panic(fmt.Errorf("[路由方式] %s 内访问了 %d 次", time.Since(v.GetLastVisitor()), v.GetVisitTimes()))
   	}
   }
   ```

   
