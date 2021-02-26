package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
)

//创建工作队列模式下的RabbitMQ实例(消息持久化，手动确认)
func NewWorkQueueRabbitMQ(url, queueName string) (*rabbitMQ, error) {
	return newRabbitMQ(url, queueName, "", "", workQueue)
}

//工作队列模式：生产者客户端
func (this *rabbitMQ) workQueueSend(msg string) (err error) {
	queue, err := this.channel.QueueDeclare(
		this.queueName,  //队列名称
		true,   //消息持久化(不持久化重启后未消费的消息会丢失)
		false, //是否自动删除(当最后一个消费者断开连接以后，是否把我们的消息从队列中删除)
		false,  //是否具有排他性
		false,    //是否阻塞
		nil,        //额外参数
	)
	if err != nil {
		return err
	}

	err = this.channel.Publish(
		this.exchange, //交换机
		queue.Name,    //routing key
		false, //如果为true，根据exchange类型和routingKey规则，如果无法找到符合条件的队列，那么就会把发送的消息返还发送者
		false, //如果为true，当exchange发送消息到队列后发现队列上没有绑定消费者，则会把消息返还给发送者
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

//工作队列模式：消费者客户端
func (this *rabbitMQ) workQueueConsume(handleFunc HandleFunc) (err error) {
	queue, err := this.channel.QueueDeclare(
		this.queueName,  //队列名称
		true,   //消息持久化
		false, //是否自动删除
		false,  //是否具有排他性
		false,    //是否阻塞
		nil,        //额外参数
	)
	if err != nil {
		return err
	}

	messages, err := this.channel.Consume(
		queue.Name, //队列名称
		"", //用来区分多个消费者
		false, //是否自动应答
		false, //是否具有排他性
		false, //如果设置为true,表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false, //队列是否阻塞
		nil,    //其它参数
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