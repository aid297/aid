package httpLimiter

import (
	"time"
)

type (
	// Visit 访问记录
	Visit struct {
		// 最后一次请求时间
		lastVisit time.Time
		// 对应Time窗口内的访问次数
		visitTimes uint16
	}

	// IpLimiter ip限流器
	IpLimiter struct{ visitMap map[string]*Visit }
)

var (
	VisitApp Visit
)

func (*Visit) New() *Visit {
	return &Visit{lastVisit: time.Now(), visitTimes: 1}
}

func (*IpLimiter) New() *IpLimiter { return NewIpLimiter() }

// NewIpLimiter 实例化：Ip 限流
//
//go:fix 推荐使用New方法
func NewIpLimiter() *IpLimiter { return &IpLimiter{visitMap: make(map[string]*Visit)} }

// Affirm 检查限流
func (my *IpLimiter) Affirm(ip string, t time.Duration, maxVisitTimes uint16) (*Visit, bool) {
	if maxVisitTimes == 0 || t == 0 {
		return nil, true
	}

	v, ok := my.visitMap[ip]
	if !ok {
		my.visitMap[ip] = VisitApp.New()
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
