package httpLimiter

import (
	"time"
)

type (
	// Visit 访问记录
	Visit struct {
		// 最后一次请求时间
		lastVisit time.Time
		// 对应 Time 窗口内的访问次数
		visitTimes uint16
	}

	// IPLimiter ip限流器
	IPLimiter struct{ visitMap map[string]*Visit }
)

func (*Visit) New() *Visit {
	return &Visit{lastVisit: time.Now(), visitTimes: 1}
}

func (*IPLimiter) New() *IPLimiter { return NewIPLimiter() }

// NewIPLimiter 实例化：Ip 限流
//
//go:fix 推荐使用New方法
func NewIPLimiter() *IPLimiter { return &IPLimiter{visitMap: make(map[string]*Visit)} }

// Affirm 检查限流
func (my *IPLimiter) Affirm(ip string, t time.Duration, maxVisitTimes uint16) (*Visit, bool) {
	if maxVisitTimes == 0 || t == 0 {
		return nil, true
	}

	v, ok := my.visitMap[ip]
	if !ok {
		my.visitMap[ip] = APP.Visit.New()
		return nil, true
	}

	if time.Since(v.lastVisit) > t {
		v.visitTimes = 1
	} else {
		v.visitTimes++
		if v.visitTimes > maxVisitTimes {
			return v, false
		}
	}
	v.lastVisit = time.Now()

	return nil, true
}

// GetLastVisitor 获取最后访问时间
func (r *Visit) GetLastVisitor() time.Time { return r.lastVisit }

// GetVisitTimes 获取窗口期内访问次数
func (r *Visit) GetVisitTimes() uint16 { return r.visitTimes }
