package mq

import (
	"bookzone/util/log"
	"github.com/streadway/amqp"
)

type RouteMessageQueue struct {
	BaseMQ
	exchange 		string
}

func NewRouteMessageQueue(exchange string) IMessageQueue {
	var err error
	routeRabbitMQ := &RouteMessageQueue{}
	routeRabbitMQ.rabbitConn, err = amqp.Dial(rabbitMQUrl)
	if err != nil {
		log.Errorf(err.Error())
		panic(err)
	}
	routeRabbitMQ.channel, err = routeRabbitMQ.rabbitConn.Channel()
	if err != nil {
		log.Errorf(err.Error())
		panic(err)
	}
	routeRabbitMQ.exchange = exchange
	return routeRabbitMQ
}

func (this *RouteMessageQueue) Publish(routingKey string, content string) error {
	var err error
	err = this.channel.ExchangeDeclare(
		this.exchange,
		"direct",
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

func (this *RouteMessageQueue) Consume(bindingKey string) (<-chan amqp.Delivery, error) {
	var err error
	err = this.channel.ExchangeDeclare(
		this.exchange,
		"direct",
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