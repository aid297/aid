package httpLimiter

import (
	"sync"
	"time"
)

type (
	// visitor 访问者对象
	visitor struct {
		ipLimiter     *IPLimiter
		t             time.Duration
		maxVisitTimes uint16
	}

	// RouteLimiter 路由限流器
	RouteLimiter struct{ RouteSetMap *sync.Map }
)

var (
	routerLimiterOnce = sync.Once{}
	routerLimiterIns  *RouteLimiter
	visitorApp        visitor
)

func (*visitor) New(t time.Duration, maxVisitTimes uint16) *visitor {
	return &visitor{ipLimiter: NewIPLimiter(), t: t, maxVisitTimes: maxVisitTimes}
}

func (*RouteLimiter) Once() *RouteLimiter { return OnceRouteLimiter() }

// OnceRouteLimiter 单例化：路由限流
//
//go:fix 推荐使用Once方法
func OnceRouteLimiter() *RouteLimiter {
	routerLimiterOnce.Do(func() { routerLimiterIns = &RouteLimiter{RouteSetMap: &sync.Map{}} })

	return routerLimiterIns
}

// Add 添加限流规则
func (my *RouteLimiter) Add(router string, t time.Duration, maxVisitTimes uint16) *RouteLimiter {
	my.RouteSetMap.Store(router, visitorApp.New(t, maxVisitTimes))

	return my
}

// Affirm 检查是否通过限流
func (my *RouteLimiter) Affirm(router, ip string) (*Visit, bool) {
	if val, exist := my.RouteSetMap.Load(router); exist {
		v := val.(*visitor)
		return v.ipLimiter.Affirm(ip, v.t, v.maxVisitTimes)
	}

	return nil, true
}
