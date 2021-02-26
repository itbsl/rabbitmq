package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
)

//创建话题模式下的RabbitMQ实例(消息持久化，手动应答)
func NewTopicRabbitMQ(url, exchange, routingKey string) (*rabbitMQ, error) {
	return newRabbitMQ(url, "", exchange, routingKey, route)
}

//话题模式：生产者客户端
func (this *rabbitMQ) topicSend(msg string) (err error) {
	//声明交换机
	err = this.channel.ExchangeDeclare(
		this.exchange,
		"topic",
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
		this.key, //要设置
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

//话题模式：消费者客户端
func (this *rabbitMQ) topicConsume(handleFunc HandleFunc) (err error) {
	//创建交换机
	err = this.channel.ExchangeDeclare(
		this.exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil
	}

	//创建队列
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

	//绑定队列到exchange中
	err = this.channel.QueueBind(
		queue.Name,
		this.key,
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
			if handleFunc(string(data.Body)) {
				data.Ack(false)
			}
		}
	}()
	fmt.Printf("[*] 等待消息，使用Ctrl+C退出\n")
	<-forever
	return nil
}



