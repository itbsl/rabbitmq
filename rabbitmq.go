package rabbitmq

import (
	"github.com/streadway/amqp"
)

/*
rabbitmq类型
simple:普通类型
workQueue:工作队列类型
pubSub:发布订阅类型
route:路由类型
topic:话题类型
*/
const (
	simple = iota
	workQueue
	pubSub
	route
	topic
)

type rabbitMQ struct {
	//连接
	conn      *amqp.Connection
	//通道
	channel   *amqp.Channel
	//队列名称
	queueName string
	//交换机
	exchange  string
	//RoutingKey/BindingKey
	key       string
	//类型
	rType     int
}

//消息处理函数
type HandleFunc func(msg string) bool

//创建RabbitMQ实例
func newRabbitMQ(url string, queueName string, exchange string, key string, rType int) (mq *rabbitMQ, err error) {
	mq = &rabbitMQ{
		queueName: queueName,
		exchange:  exchange,
		key:       key,
		rType:     rType,
	}

	//创建连接
	mq.conn, err = amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	//创建通道
	mq.channel, err = mq.conn.Channel()
	if err != nil {
		return nil, err
	}

	return mq, nil
}

//断开conn和channel
func (this *rabbitMQ) Close() {
	this.channel.Close()
	this.conn.Close()
}

//生产者客户端接口
func (this *rabbitMQ) Send(msg string) error {
	switch this.rType {
	case simple:
		return this.simpleSend(msg)
	case workQueue:
		return this.workQueueSend(msg)
	case pubSub:
		return this.pubSubSend(msg)
	case route:
		return this.routeSend(msg)
	case topic:
		return this.topicSend(msg)
	}
	return nil
}

//消费者客户端接口
func (this *rabbitMQ) Consume(handleFunc HandleFunc) error {
	switch this.rType {
	case simple:
		return this.simpleConsume(handleFunc)
	case workQueue:
		return this.workQueueConsume(handleFunc)
	case pubSub:
		return this.pubSubConsume(handleFunc)
	case route:
		return this.routeConsume(handleFunc)
	case topic:
		return this.topicConsume(handleFunc)
	}
	return nil
}