package websocketPool

import "time"

type (
	// Heart 链接心跳
	Heart struct {
		ticker *time.Ticker
		fn     func(*Client)
	}

	// MessageTimeout 通信超时
	MessageTimeout struct{ interval time.Duration }
)

var HeartApp Heart

func (*Heart) New() *Heart { return &Heart{} }

// NewHeart 实例化：链接心跳
//
//go:fix 推荐使用：New方法
func NewHeart() *Heart { return &Heart{} }

// SetInterval 设置定时器
func (my *Heart) SetInterval(interval time.Duration) *Heart {
	if my.ticker != nil {
		my.ticker.Reset(interval)
	} else {
		my.ticker = time.NewTicker(interval)
	}

	return my
}

// SetFn 设置回调：定时器执行内容
func (my *Heart) SetFn(fn func(client *Client)) *Heart {
	my.fn = fn

	return my
}

// Stop 停止定时器
func (my *Heart) Stop() *Heart {
	my.ticker.Stop()

	return my
}

// DefaultHeart 默认心跳：10秒
func DefaultHeart() *Heart {
	// return NewHeart().SetInterval(time.Second * 10).SetFn(func(client *Client) {
	// _, _ = client.SendMsg(MsgType.Ping(), []byte("ping"))
	// })
	return NewHeart().SetInterval(60 * time.Second).SetFn(nil)
}

// NewMessageTimeout 实例化：链接超时
func NewMessageTimeout() *MessageTimeout { return &MessageTimeout{} }

// SetInterval 设置定时器时间
func (r *MessageTimeout) SetInterval(interval time.Duration) *MessageTimeout {
	r.interval = interval

	return r
}

// DefaultMessageTimeout 默认消息超时：5秒
func DefaultMessageTimeout() *MessageTimeout { return NewMessageTimeout().SetInterval(time.Second * 5) }
