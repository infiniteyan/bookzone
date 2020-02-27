package mq

import (
	"bookzone/util/log"
	"github.com/streadway/amqp"
)

type TopicMessageQueue struct {
	BaseMQ
	exchange 		string
}

func NewTopicMessageQueue(exchange string) IMessageQueue {
	var err error
	topicRabbitMQ := &TopicMessageQueue{}
	topicRabbitMQ.rabbitConn, err = amqp.Dial(rabbitMQUrl)
	if err != nil {
		log.Errorf(err.Error())
		panic(err)
	}
	topicRabbitMQ.channel, err = topicRabbitMQ.rabbitConn.Channel()
	if err != nil {
		log.Errorf(err.Error())
		panic(err)
	}
	topicRabbitMQ.exchange = exchange
	return topicRabbitMQ
}

func (this *TopicMessageQueue) Publish(routingKey string, content string) error {
	var err error
	err = this.channel.ExchangeDeclare(
		this.exchange,
		"topic",
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
		routingKey,
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

func (this *TopicMessageQueue) Consume(bindingKey string) (<-chan amqp.Delivery, error) {
	var err error
	err = this.channel.ExchangeDeclare(
		this.exchange,
		"topic",
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
		bindingKey,
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