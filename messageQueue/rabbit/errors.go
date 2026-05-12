package rabbit

import (
	"reflect"

	"github.com/aid297/aid/v2/anySlice"
	"github.com/aid297/aid/v2/myError"
)

type (
	ConnRabbitError       struct{ myError.MyError }
	NewChannelError       struct{ myError.MyError }
	NewQueueError         struct{ myError.MyError }
	QueueNotExistError    struct{ myError.MyError }
	PublishMessageError   struct{ myError.MyError }
	RegisterConsumerError struct{ myError.MyError }
)

var (
	ConnRabbitErr       ConnRabbitError
	NewChannelErr       NewChannelError
	NewQueueErr         NewQueueError
	QueueNotExistErr    QueueNotExistError
	PublishMessageErr   PublishMessageError
	RegisterConsumerErr RegisterConsumerError
)

func (*ConnRabbitError) New(msg string) myError.IMyError {
	return &ConnRabbitError{myError.MyError{Msg: anySlice.New(anySlice.Items("链接rabbit-mq错误", msg)).JoinNotEmpty("：")}}
}
func (*ConnRabbitError) Wrap(err error) myError.IMyError {
	return &ConnRabbitError{myError.MyError{Msg: anySlice.New(anySlice.Items("链接rabbit-mq错误", err.Error())).JoinNotEmpty("：")}}
}

func (*ConnRabbitError) Panic() myError.IMyError {
	return &ConnRabbitError{myError.MyError{Msg: "链接rabbit-mq错误"}}
}

func (my *ConnRabbitError) Error() string { return my.Msg }

func (my *ConnRabbitError) Is(target error) bool { return reflect.DeepEqual(target, my) }

func (*NewChannelError) New(msg string) myError.IMyError {
	return &NewChannelError{myError.MyError{Msg: anySlice.New(anySlice.Items("创建channel错误", msg)).JoinNotEmpty("：")}}
}
func (*NewChannelError) Wrap(err error) myError.IMyError {
	return &NewChannelError{myError.MyError{Msg: anySlice.New(anySlice.Items("创建channel错误", err.Error())).JoinNotEmpty("：")}}
}

func (*NewChannelError) Panic() myError.IMyError {
	return &NewChannelError{myError.MyError{Msg: "创建channel错误"}}
}

func (my *NewChannelError) Error() string { return my.Msg }

func (my *NewChannelError) Is(target error) bool { return reflect.DeepEqual(target, my) }

func (*NewQueueError) New(msg string) myError.IMyError {
	return &NewQueueError{myError.MyError{Msg: anySlice.New(anySlice.Items("创建队列错误", msg)).JoinNotEmpty("：")}}
}
func (*NewQueueError) Wrap(err error) myError.IMyError {
	return &NewQueueError{myError.MyError{Msg: anySlice.New(anySlice.Items("创建队列错误", err.Error())).JoinNotEmpty("：")}}
}

func (*NewQueueError) Panic() myError.IMyError {
	return &NewQueueError{myError.MyError{Msg: "创建队列错误"}}
}

func (my *NewQueueError) Error() string { return my.Msg }

func (my *NewQueueError) Is(target error) bool { return reflect.DeepEqual(target, my) }

func (*QueueNotExistError) New(msg string) myError.IMyError {
	return &QueueNotExistError{myError.MyError{Msg: anySlice.New(anySlice.Items("队列不存在", msg)).JoinNotEmpty("：")}}
}
func (*QueueNotExistError) Wrap(err error) myError.IMyError {
	return &QueueNotExistError{myError.MyError{Msg: anySlice.New(anySlice.Items("队列不存在", err.Error())).JoinNotEmpty("：")}}
}

func (*QueueNotExistError) Panic() myError.IMyError {
	return &QueueNotExistError{myError.MyError{Msg: "队列不存在"}}
}

func (my *QueueNotExistError) Error() string { return my.Msg }

func (my *QueueNotExistError) Is(target error) bool { return reflect.DeepEqual(target, my) }

func (*PublishMessageError) New(msg string) myError.IMyError {
	return &PublishMessageError{myError.MyError{Msg: anySlice.New(anySlice.Items("生产消息错误", msg)).JoinNotEmpty("：")}}
}
func (*PublishMessageError) Wrap(err error) myError.IMyError {
	return &PublishMessageError{myError.MyError{Msg: anySlice.New(anySlice.Items("生产消息错误", err.Error())).JoinNotEmpty("：")}}
}

func (*PublishMessageError) Panic() myError.IMyError {
	return &PublishMessageError{myError.MyError{Msg: "生产消息错误"}}
}

func (my *PublishMessageError) Error() string { return my.Msg }

func (my *PublishMessageError) Is(target error) bool { return reflect.DeepEqual(target, my) }

func (*RegisterConsumerError) New(msg string) myError.IMyError {
	return &RegisterConsumerError{myError.MyError{Msg: anySlice.New(anySlice.Items("注册消费者错误", msg)).JoinNotEmpty("：")}}
}
func (*RegisterConsumerError) Wrap(err error) myError.IMyError {
	return &RegisterConsumerError{myError.MyError{Msg: anySlice.New(anySlice.Items("注册消费者错误", err.Error())).JoinNotEmpty("：")}}
}

func (*RegisterConsumerError) Panic() myError.IMyError {
	return &RegisterConsumerError{myError.MyError{Msg: "注册消费者错误"}}
}

func (my *RegisterConsumerError) Error() string { return my.Msg }

func (my *RegisterConsumerError) Is(target error) bool { return reflect.DeepEqual(target, my) }
