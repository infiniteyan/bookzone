package mq

import (
	"bookzone/util/log"
	"github.com/streadway/amqp"
)

type PubSubMessageQueue struct {
	BaseMQ
	exchange 		string
}

func NewPubSubMessageQueue(exchange string) IMessageQueue {
	var err error
	pubsubRabbitMQ := &PubSubMessageQueue{}
	pubsubRabbitMQ.rabbitConn, err = amqp.Dial(rabbitMQUrl)
	if err != nil {
		log.Errorf(err.Error())
		panic(err)
	}
	pubsubRabbitMQ.channel, err = pubsubRabbitMQ.rabbitConn.Channel()
	if err != nil {
		log.Errorf(err.Error())
		panic(err)
	}
	pubsubRabbitMQ.exchange = exchange
	return pubsubRabbitMQ
}

func (this *PubSubMessageQueue) Publish(routingKey string, content string) error {
	var err error
	err = this.channel.ExchangeDeclare(
		this.exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Errorf(err.Error())
		return err
	}

	err = this.channel.Publish(
		this.exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType:"text/plain",
			Body: []byte(content),
		})
	if err != nil {
		log.Errorf(err.Error())
		return err
	}

	return nil
}

func (this *PubSubMessageQueue) Consume(bindingKey string) (<-chan amqp.Delivery, error) {
	var err error
	err = this.channel.ExchangeDeclare(
		this.exchange,
		"fanout",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}

	queue, err := this.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil)
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}

	err = this.channel.QueueBind(
		queue.Name,
		"",
		this.exchange,
		false,
		nil)
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}

	msgs, err := this.channel.Consume(
		queue.Name,
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