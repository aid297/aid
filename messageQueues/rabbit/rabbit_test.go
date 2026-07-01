package rabbit

import (
	"testing"
	"time"
)

func Test1(t *testing.T) {
	rabbit := APP.Rabbit.New("admin", "jcyf@cbit", "127.0.0.1", "5672", "")
	defer func() { _ = rabbit.Close() }()

	pool := APP.Pool.Once().Set("default", rabbit)

	rabbit = pool.Get("default")
	if rabbit == nil {
		t.Fatalf("没有找到链接：%s", "default")
	}

	rabbit.NewQueue("message")
	if rabbit.Error() != nil {
		t.Fatalf("创建队列失败：%v", rabbit.Error())
	}

	rabbit.Publish("message", "hello world"+time.Now().Format(time.DateTime))
}

func Test2(t *testing.T) {
	rabbit := APP.Rabbit.New("admin", "jcyf@cbit", "127.0.0.1", "5672", "")
	defer func() { _ = rabbit.Close() }()

	rabbit.NewQueue("message")
	consumer := rabbit.Consume("message", "", func(prototypeMessage []byte) error {
		message := APP.Message.Parse(prototypeMessage)
		t.Logf("收到消息：%s", message.Content)
		return nil
	})
	go consumer.Go()
}
