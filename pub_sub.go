package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
)

//创建发布/订阅模式下的RabbitMQ实例(消息持久化，自动应答)
func NewPubSubRabbitMQ(url, exchange string) (*rabbitMQ, error) {
	return newRabbitMQ(url, "", exchange, "", pubSub)
}

//发布订阅模式：生产者客户端
func (this *rabbitMQ) pubSubSend(msg string) (err error) {
	//1.声明交换机
	err = this.channel.ExchangeDeclare(
		this.exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	//发送消息
	err = this.channel.Publish(
		this.exchange,
		"",
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	if err != nil {
		return err
	}

	return nil
}

//发布订阅模式：消费者客户端
func (this *rabbitMQ) pubSubConsume(handleFunc HandleFunc) (err error) {
	//声明交换机
	err = this.channel.ExchangeDeclare(
		this.exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	//声明队列
	queue, err := this.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	//队列绑定
	err = this.channel.QueueBind(
		queue.Name,
		"", //在pub/sub模式下，这里的key要为空
		this.exchange,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	//注册消费者
	messages, err := this.channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for data := range messages {
			//实现我们要处理的逻辑
			handleFunc(string(data.Body))
		}
	}()
	fmt.Printf("[*] 等待消息，使用Ctrl+C退出\n")
	<-forever
	return nil
}

