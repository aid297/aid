package websocketPool

import "time"

type (
	Heart struct {
		ticker *time.Ticker
		fn     func(c *Client)
	}

	HeartAttributer interface{ Register(heart *Heart) }

	AttrHeartTimeout struct{ interval time.Duration }
	AttrHeartFn      struct{ fn func(c *Client) }
)

// DefaultHeart 默认心跳：10秒
func DefaultHeart() *Heart {
	// return NewHeart().SetInterval(time.Second * 10).SetFn(func(client *Client) {
	// _, _ = client.SendMsg(MsgType.Ping(), []byte("ping"))
	// })
	return APP.Heart.New(HeartInterval(10*time.Second), HeartFn(func(c *Client) { c.SendMsg(MsgType.Ping(), nil) }))
}

func (*Heart) New(attrs ...HeartAttributer) *Heart {
	ins := Heart{}
	for _, attr := range attrs {
		attr.Register(&ins)
	}
	return &ins
}

func HeartInterval(interval time.Duration) HeartAttributer {
	return AttrHeartTimeout{interval: interval}
}
func (my AttrHeartTimeout) Register(heart *Heart) {
	heart.ticker = time.NewTicker(my.interval)
}

func HeartFn(fn func(c *Client)) HeartAttributer { return AttrHeartFn{fn: fn} }
func (my AttrHeartFn) Register(heart *Heart)     { heart.fn = my.fn }
