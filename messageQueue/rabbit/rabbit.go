package rabbit

import (
	"fmt"
	"sync"

	"github.com/streadway/amqp"
)

type (
	Rabbit struct {
		err         error
		username    string
		password    string
		host        string
		port        string
		virtualHost string
		conn        *amqp.Connection
		ch          *amqp.Channel
		queues      map[string]amqp.Queue
		mu          sync.RWMutex
	}
)

// New 创建一个 Rabbit 实例
func (*Rabbit) New(
	username,
	password,
	host,
	port,
	virtualHost string,
) *Rabbit {
	ins := &Rabbit{
		username:    username,
		password:    password,
		host:        host,
		port:        port,
		virtualHost: virtualHost,
		queues:      make(map[string]amqp.Queue),
	}

	// 连接到 RabbitMQ
	ins.conn, ins.err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/%s", ins.username, ins.password, ins.host, ins.port, ins.virtualHost))
	if ins.err != nil {
		ins.err = ConnRabbitErr.Wrap(ins.err)
	}

	ins.NewChannel() // 创建频道

	return ins
}

// 获取链接
func (my *Rabbit) getConn() *amqp.Connection { return my.conn }

// GetConn 获取链接
func (my *Rabbit) GetConn() *amqp.Connection {
	my.mu.RLock()
	defer my.mu.RLock()

	return my.getConn()
}

// Error 获取错误
func (my *Rabbit) Error() error { return my.err }

// Close 关闭链接
func (my *Rabbit) Close() error {
	my.mu.Lock()
	defer my.mu.Unlock()

	if my.conn != nil {
		my.closeChannel()
		return my.conn.Close()
	}

	return nil
}

// newChannel 创建频道
func (my *Rabbit) newChannel() {
	if my.ch == nil {
		my.ch, my.err = my.getConn().Channel()
	}
}

// NewChannel 创建频道
func (my *Rabbit) NewChannel() *Rabbit {
	my.mu.Lock()
	defer my.mu.Unlock()

	my.newChannel()

	return my
}

// closeChannel 关闭频道
func (my *Rabbit) closeChannel() {
	if my.ch != nil {
		my.err = my.ch.Close()
	}
}

// CloseChannel 关闭频道
func (my *Rabbit) CloseChannel() *Rabbit {
	my.mu.Lock()
	defer my.mu.Unlock()

	my.closeChannel()

	return my
}

// newQueue 创建队列
func (my *Rabbit) newQueue(queueName string) amqp.Queue {
	var queue amqp.Queue
	queue, my.err = my.ch.QueueDeclare(
		queueName, // 队列名称
		true,      // 持久化
		false,     // 自动删除
		false,     // 独占
		false,     // 不等待
		nil,       // 附加属性
	)
	if my.err != nil {
		my.err = NewQueueErr.Wrap(my.err)
	}
	return queue
}

// NewQueue 创建队列
func (my *Rabbit) NewQueue(queueName string) *Rabbit {
	my.mu.Lock()
	defer my.mu.Unlock()

	if my.err != nil {
		return my
	}

	my.newChannel()
	if my.err != nil {
		return my
	}

	my.queues[queueName] = my.newQueue(queueName)

	return my
}

// getQueue 获取队列
func (my *Rabbit) getQueue(queueName string) amqp.Queue {
	if queue, ok := my.queues[queueName]; ok {
		return queue
	} else {
		my.err = QueueNotExistErr.New(queueName)
	}

	return amqp.Queue{}
}

// Publish 生产消息
func (my *Rabbit) Publish(queueName string, body string) *Rabbit {
	my.mu.Lock()
	defer my.mu.Unlock()

	if my.err != nil {
		return my
	}

	queue := my.getQueue(queueName)
	if my.err != nil {
		return my
	}

	// 发送消息
	my.err = my.ch.Publish(
		"",         // 默认交换机
		queue.Name, // 路由键，使用队列名称
		false,      // 是否立即发送
		false,      // 是否持久化
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if my.err != nil {
		my.err = PublishMessageErr.Wrap(my.err)
	}

	return my
}

// Consume 消费消息
func (my *Rabbit) Consume(queueName, consumer string, parseFn func(prototypeMessage []byte) error) *Consumer {
	my.mu.RLock()
	defer my.mu.RUnlock()

	if my.ch == nil {
		return nil
	}

	my.newChannel()
	if my.err != nil {
		return nil
	}

	queue := my.getQueue(queueName)
	if my.err != nil {
		return nil
	}

	return APP.Consumer.New(my.ch, queue, consumer, parseFn)
}
