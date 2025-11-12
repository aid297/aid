package rabbit

import (
	"fmt"
	"reflect"

	"github.com/aid297/aid/array"
	"github.com/aid297/aid/myError"
	"github.com/aid297/aid/operation"
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
	return &ConnRabbitError{myError.MyError{Msg: array.NewDestruction("链接rabbit-mq错误", msg).JoinWithoutEmpty()}}
}
func (*ConnRabbitError) Wrap(err error) myError.IMyError {
	return &ConnRabbitError{myError.MyError{Msg: fmt.Errorf("链接rabbit-mq错误" + operation.Ternary(err != nil, "：%w", "%w")).Error()}}
}

func (*ConnRabbitError) Panic() myError.IMyError {
	return &ConnRabbitError{myError.MyError{Msg: "链接rabbit-mq错误"}}
}

func (my *ConnRabbitError) Error() string { return my.Msg }

func (my *ConnRabbitError) Is(target error) bool { return reflect.DeepEqual(target, my) }

func (*NewChannelError) New(msg string) myError.IMyError {
	return &NewChannelError{myError.MyError{Msg: array.NewDestruction("创建channel错误", msg).JoinWithoutEmpty()}}
}
func (*NewChannelError) Wrap(err error) myError.IMyError {
	return &NewChannelError{myError.MyError{Msg: fmt.Errorf("创建channel错误" + operation.Ternary(err != nil, "：%w", "%w")).Error()}}
}

func (*NewChannelError) Panic() myError.IMyError {
	return &NewChannelError{myError.MyError{Msg: "创建channel错误"}}
}

func (my *NewChannelError) Error() string { return my.Msg }

func (my *NewChannelError) Is(target error) bool { return reflect.DeepEqual(target, my) }

func (*NewQueueError) New(msg string) myError.IMyError {
	return &NewQueueError{myError.MyError{Msg: array.NewDestruction("创建队列错误", msg).JoinWithoutEmpty()}}
}
func (*NewQueueError) Wrap(err error) myError.IMyError {
	return &NewQueueError{myError.MyError{Msg: fmt.Errorf("创建队列错误" + operation.Ternary(err != nil, "：%w", "%w")).Error()}}
}

func (*NewQueueError) Panic() myError.IMyError {
	return &NewQueueError{myError.MyError{Msg: "创建队列错误"}}
}

func (my *NewQueueError) Error() string { return my.Msg }

func (my *NewQueueError) Is(target error) bool { return reflect.DeepEqual(target, my) }

func (*QueueNotExistError) New(msg string) myError.IMyError {
	return &QueueNotExistError{myError.MyError{Msg: array.NewDestruction("队列不存在", msg).JoinWithoutEmpty()}}
}
func (*QueueNotExistError) Wrap(err error) myError.IMyError {
	return &QueueNotExistError{myError.MyError{Msg: fmt.Errorf("队列不存在" + operation.Ternary(err != nil, "：%w", "%w")).Error()}}
}

func (*QueueNotExistError) Panic() myError.IMyError {
	return &QueueNotExistError{myError.MyError{Msg: "队列不存在"}}
}

func (my *QueueNotExistError) Error() string { return my.Msg }

func (my *QueueNotExistError) Is(target error) bool { return reflect.DeepEqual(target, my) }

func (*PublishMessageError) New(msg string) myError.IMyError {
	return &PublishMessageError{myError.MyError{Msg: array.NewDestruction("生产消息错误", msg).JoinWithoutEmpty()}}
}
func (*PublishMessageError) Wrap(err error) myError.IMyError {
	return &PublishMessageError{myError.MyError{Msg: fmt.Errorf("生产消息错误" + operation.Ternary(err != nil, "：%w", "%w")).Error()}}
}

func (*PublishMessageError) Panic() myError.IMyError {
	return &PublishMessageError{myError.MyError{Msg: "生产消息错误"}}
}

func (my *PublishMessageError) Error() string { return my.Msg }

func (my *PublishMessageError) Is(target error) bool { return reflect.DeepEqual(target, my) }

func (*RegisterConsumerError) New(msg string) myError.IMyError {
	return &RegisterConsumerError{myError.MyError{Msg: array.NewDestruction("注册消费者错误", msg).JoinWithoutEmpty()}}
}
func (*RegisterConsumerError) Wrap(err error) myError.IMyError {
	return &RegisterConsumerError{myError.MyError{Msg: fmt.Errorf("注册消费者错误" + operation.Ternary(err != nil, "：%w", "%w")).Error()}}
}

func (*RegisterConsumerError) Panic() myError.IMyError {
	return &RegisterConsumerError{myError.MyError{Msg: "注册消费者错误"}}
}

func (my *RegisterConsumerError) Error() string { return my.Msg }

func (my *RegisterConsumerError) Is(target error) bool { return reflect.DeepEqual(target, my) }
