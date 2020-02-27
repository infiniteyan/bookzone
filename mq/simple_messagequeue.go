package mq

import (
	"bookzone/util/log"
	"github.com/streadway/amqp"
)

type SimpleMessageQueue struct {
	BaseMQ
	queueName 		string
}

func NewSimpleMessageQueue(queueName string) IMessageQueue {
	var err error
	simpleRabbitMQ := &SimpleMessageQueue{}
	simpleRabbitMQ.rabbitConn, err = amqp.Dial(rabbitMQUrl)
	if err != nil {
		log.Errorf(err.Error())
		panic(err)
	}
	simpleRabbitMQ.channel, err = simpleRabbitMQ.rabbitConn.Channel()
	if err != nil {
		log.Errorf(err.Error())
		panic(err)
	}
	simpleRabbitMQ.queueName = queueName
	return simpleRabbitMQ
}

func (this *SimpleMessageQueue) Publish(routingKey string, content string) error {
	_, err := this.channel.QueueDeclare(
		routingKey,
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Errorf(err.Error())
		return err
	}

	err = this.channel.Publish(
		"",
		this.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(content),
		})
	return err
}
func (this *SimpleMessageQueue) Consume(routingKey string) (<-chan amqp.Delivery, error) {
	q ,err := this.channel.QueueDeclare(this.queueName,
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}

	msgs, err := this.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}
	return msgs, nil
}