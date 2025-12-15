package rabbit

import (
	"github.com/streadway/amqp"
)

type (
	Consumer struct {
		err              error
		stop             chan struct{}
		ch               *amqp.Channel
		queue            amqp.Queue
		consumer         string
		parseFn          func(message []byte) error
		prototypeMessage <-chan amqp.Delivery
		isListening      bool
	}
)

func (*Consumer) New(
	ch *amqp.Channel,
	queue amqp.Queue,
	consumer string,
	parseFn func(message []byte) error,
) *Consumer {
	return &Consumer{
		ch:               ch,
		queue:            queue,
		consumer:         consumer,
		stop:             make(chan struct{}),
		parseFn:          parseFn,
		prototypeMessage: make(chan amqp.Delivery),
	}
}

// Error 获取错误信息
func (my *Consumer) Error() error { return my.err }

func (my *Consumer) Go() <-chan amqp.Delivery {
	// 获取消息
	msgs, err := my.ch.Consume(
		my.queue.Name, // 队列名称
		my.consumer,   // 消费者名称
		true,          // 自动确认消息
		false,         // 独占
		false,         // 不等待
		false,         // 不阻塞
		nil,           // 附加属性
	)
	if err != nil {
		my.err = RegisterConsumerErr.Wrap(err)
		return nil
	}
	return msgs
}

// Start 监听：开始
func (my *Consumer) Start() *Consumer {
	defer func() { my.isListening = true }()

	msgs := my.Go()
	go func() {
		for {
			select {
			case <-my.stop:
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}
				my.parseFn(msg.Body)
			}
		}
	}()

	return my
}

// Stop 停止监听
func (my *Consumer) Stop() *Consumer {
	defer func() { my.isListening = false }()

	my.stop <- struct{}{}
	return my
}
