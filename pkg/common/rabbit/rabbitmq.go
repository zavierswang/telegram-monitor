package mq

import (
	"github.com/streadway/amqp"
	"telegram-monitor/pkg/core/logger"
)

func NewRabbitMQ(exchange, dsn, route, queue string) (MessageQueue, error) {
	var messageQueue = MessageQueue{
		ExchangeName: exchange,
		RouteKey:     route,
		QueueName:    queue,
	}

	// 建立amqp链接
	//logger.Info("rabbitMQ dsn: %s", dsn)
	conn, err := amqp.Dial(dsn)
	if err != nil {
		logger.Error("Failed to connect to RabbitMQ, %v", err)
		return messageQueue, err
	}
	messageQueue.conn = conn
	// 建立channel通道
	ch, err := conn.Channel()
	if err != nil {
		logger.Error("open channel failed %v", err)
		return messageQueue, err
	}
	messageQueue.ch = ch
	// 声明exchange交换器
	err = messageQueue.declareExchange(exchange, nil)
	if err != nil {
		logger.Error("declare exchange failed %v", err)
		return messageQueue, err
	}

	return messageQueue, nil
}
