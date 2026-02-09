### Rabbit MQ 消息队列使用说明

1. 创建`生产者`
   ```go
   package main
   
   import (
   	`errors`
   	`fmt`
   	`time`
   
   	`github.com/aid297/aid/messageQueue/rabbit`
   )
   
   func main() {
   	rbt := rabbit.APP.Rabbit.New("admin", "jcyf@cbit", "127.0.0.1", "5672", "")
   	defer func() { _ = rbt.Close() }()
   
   	pool := rabbit.APP.Pool.Once().Set("default", rbt)
   
   	rbt = pool.Get("default")
   	if rbt == nil {
   		panic(errors.New("没有找到链接：default"))
   	}
   
   	rbt.NewQueue("message")
   	if rbt.Error() != nil {
   		panic(fmt.Errorf("创建队列失败：%v", rbt.Error()))
   	}
   
   	rbt.Publish("message", "hello world"+time.Now().Format(time.DateTime))
   }
   ```

2. 创建`消费者`
   ```go
   package main
   
   import (
   	`log`
   
   	`github.com/aid297/aid/messageQueue/rabbit`
   )
   
   func main() {
   	rbt := rabbit.APP.Rabbit.New("admin", "jcyf@cbit", "127.0.0.1", "5672", "")
   	defer func() { _ = rbt.Close() }()
   
   	rbt.NewQueue("message")
   	consumer := rbt.Consume("message", "", func(prototypeMessage []byte) error {
   		message := rabbit.APP.Message.Parse(prototypeMessage)
   		log.Printf("接收消息错误：%v", message.Content)
   		return nil
   	})
   	go consumer.Go()
   }
   ```

   