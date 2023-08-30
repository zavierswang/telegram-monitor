package mq

import (
	"github.com/streadway/amqp"
	"strconv"
	"telegram-monitor/pkg/core/logger"
)

type Message struct {
	DelayTime int
	Body      []byte
}

type MessageQueue struct {
	conn         *amqp.Connection // amqp链接对象
	ch           *amqp.Channel    // channel对象
	ExchangeName string           // 交换器名称
	RouteKey     string           // 路由名称
	QueueName    string           // 队列名称
}

// Consumer 消费者回调方法
type Consumer func(amqp.Delivery)

// SendMessage 发送普通消息
func (mq *MessageQueue) SendMessage(message Message) error {
	err := mq.ch.Publish(
		mq.ExchangeName, // exchange
		mq.RouteKey,     // route key
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message.Body,
		},
	)
	return err
}

// SendDelayMessage 发送延迟消息
func (mq *MessageQueue) SendDelayMessage(message Message) error {
	delayQueueName := mq.QueueName + "_delay:" + strconv.Itoa(message.DelayTime)
	delayRouteKey := mq.RouteKey + "_delay:" + strconv.Itoa(message.DelayTime)

	// 定义延迟队列(死信队列)
	dq, err := mq.declareQueue(
		delayQueueName,
		amqp.Table{
			"x-dead-letter-exchange":    mq.ExchangeName, // 指定死信交换机
			"x-dead-letter-routing-key": mq.RouteKey,     // 指定死信routing-key
		},
	)
	if err != nil {
		return err
	}

	// 延迟队列绑定到exchange
	mq.bindQueue(dq.Name, delayRouteKey, mq.ExchangeName)

	// 发送消息，将消息发送到延迟队列，到期后自动路由到正常队列中
	err = mq.ch.Publish(
		mq.ExchangeName,
		delayRouteKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message.Body),
			Expiration:  strconv.Itoa(message.DelayTime * 1000),
		},
	)
	return err
}

// Consume 获取消费消息
func (mq *MessageQueue) Consume(fn Consumer) error {
	// 声明队列
	q, err := mq.declareQueue(mq.QueueName, nil)
	if err != nil {
		return err
	}

	// 队列绑定到exchange
	mq.bindQueue(q.Name, mq.RouteKey, mq.ExchangeName)

	// 设置Qos
	err = mq.ch.Qos(1, 0, false)
	if err != nil {
		return err
	}

	// 监听消息
	msg, err := mq.ch.Consume(
		q.Name, // queue name,
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	// forever := make(chan bool), 注册在主进程，不需要阻塞
	go func() {
		for d := range msg {
			fn(d)
			_ = d.Ack(false)
		}
	}()

	logger.Info(" [*] Waiting for logs. To exit press CTRL+C")
	return nil
}

// Close 关闭链接
func (mq *MessageQueue) Close() {
	mq.ch.Close()
	mq.conn.Close()
}

// declareQueue 定义队列
func (mq *MessageQueue) declareQueue(name string, args amqp.Table) (amqp.Queue, error) {
	q, err := mq.ch.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		args,
	)
	if err != nil {
		return q, err
	}

	return q, nil
}

func (mq *MessageQueue) declareExchange(exchange string, args amqp.Table) error {
	err := mq.ch.ExchangeDeclare(
		exchange,
		"direct",
		true,
		false,
		false,
		false,
		args,
	)
	return err
}

// bindQueue 绑定队列
func (mq *MessageQueue) bindQueue(queue, route, exchange string) error {
	err := mq.ch.QueueBind(
		queue,
		route,
		exchange,
		false,
		nil,
	)
	return err
}
