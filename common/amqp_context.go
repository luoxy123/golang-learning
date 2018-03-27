package common

import (
	"fmt"
	"time"

	logstash "github.com/FeiniuBus/log"
	"github.com/streadway/amqp"
)

type AmqpContext struct {
	queryUrl  string
	queueName string
	logger    *logstash.Logger
	Data      chan []byte
}

func NewAmqpContext(o RabbitOptions, queueName string, logger *logstash.Logger) *AmqpContext {
	queryUrl := fmt.Sprintf("amqp://%s:%s@%s:%d/", o.UserName, o.Password, o.Host, o.Port)
	context := &AmqpContext{
		queryUrl:  queryUrl,
		queueName: queueName,
		logger:    logger,
	}
	context.Data = make(chan []byte, 1000)
	return context
}
func (c *AmqpContext) Connect() {
start:
	errorChan := make(chan *amqp.Error)
	conn, err := amqp.Dial(c.queryUrl)
	for err != nil {
		c.logger.Error(err.Error())
		time.Sleep(3 * time.Second)
		conn, err = amqp.Dial(c.queryUrl)
	}
	conn.NotifyClose(errorChan)
	channel, err := conn.Channel()
	if err != nil {
		c.logger.Error(err.Error())
		return
	}
	queue, err := channel.QueueDeclare(c.queueName, true, false, false, false, nil)
	if err != nil {
		c.logger.Error(err.Error())
		return
	}
	deliveries, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		c.logger.Error(err.Error())
		return
	}
	for {
		select {
		case msg := <-deliveries:
			c.Data <- msg.Body
			c.logger.Info("收到队列消息:", string(msg.Body))
		case err := <-errorChan:
			if err != nil {
				c.logger.Error("rabbitmq 连接断开，即将重连。错误信息：", *err)
				goto start
			}
		}
	}
}
