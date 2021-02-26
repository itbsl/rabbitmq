## 安装

```shell
go get github.com/itbsl/rabbitmq
```

## 使用

使用起来，rabbitmq只包含5个创建不同类型的RabbitMQ实例的方法，和一个关闭连接的方法，以及发送消息和获取消息的的两个方法，共八个方法
方法如下：

```go
//创建普通类型的rabbitmq实例，消息没有持久化
rabbitmq.NewSimpleRabbitMQ()
//创建工作队列模式的RabbitMQ实例：消息持久化，需要手动应答
rabbitmq.NewWorkQueueRabbitMQ()
//创建发布/订阅模式的RabbitMQ实例：消息持久化，自动应答
rabbitmq.NewPubSubRabbitMQ()
//创建路由模式的RabbitMQ实例：消息持久化，手动应答
rabbitmq.NewRouteRabbitMQ()
//创建话题模式的RabbitMQ实例：消息持久化，手动应答
rabbitmq.NewTopicRabbitMQ()
```

示例：

生产者客户端
1.创建RabbitMQ实例
2.发送消息
```go
package main

import (
	"github.com/itbsl/rabbitmq"
	"log"
)

func main() {
	//连接RabbitMQ
	url := "amqp://root:root@127.0.0.1:5672/"
	mq, err := rabbitmq.NewWorkQueueRabbitMQ(url, "hello")
	if err != nil {
		log.Fatalf("连接RabbitMQ出错，错误为：%v\n", err)
	}
	defer mq.Close()

	//发送消息
	err = mq.Send("hello world!")
	if err != nil {
		log.Fatalf("发送消息失败，错误为：%v\n", err)
	}
}
```

消费者客户端
1.创建所需类型的RabbitMQ实例
2.消费消息，业务代码写在回调函数里
```go
package main

import (
	"github.com/itbsl/rabbitmq"
	"log"
)

func main() {
	//连接RabbitMQ
	url := "amqp://root:root@127.0.0.1:5672/"
	mq, err := rabbitmq.NewWorkQueueRabbitMQ(url, "hello")
	if err != nil {
		log.Fatalf("连接RabbitMQ出错，错误为：%v\n", err)
	}
	defer mq.Close()

	//消费消息
	err = mq.Consume(func(msg string) bool {
		//这里写业务处理代码，业务处理成功，返回true用于手动应答，失败返回false,消息会被重新放入队列
		return true
	})
	if err != nil {
		log.Fatalf("消费失败，错误为：%v\n", err)
	}
}
```