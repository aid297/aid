package websocketPool

import "time"

type (
	// MessageTimeout 通信超时
	MessageTimeout struct{ interval time.Duration }

	MessageTimeoutAttributer interface{ Register(timeout *MessageTimeout) }

	AttrMessageTimeoutInterval struct{ interval time.Duration }
)

// DefaultMessageTimeout 默认消息超时：5秒
func DefaultMessageTimeout() *MessageTimeout {
	return APP.MessageTimeout.New(MessageTimeoutInterval(5 * time.Second))
}

func (*MessageTimeout) New(attrs ...MessageTimeoutAttributer) *MessageTimeout {
	ins := MessageTimeout{}
	for _, attr := range attrs {
		attr.Register(&ins)
	}
	return &ins
}

func MessageTimeoutInterval(interval time.Duration) MessageTimeoutAttributer {
	return AttrMessageTimeoutInterval{interval: interval}
}
func (my AttrMessageTimeoutInterval) Register(timeout *MessageTimeout) {
	timeout.interval = my.interval
}
