package mq

import (
	"bookzone/util/log"
	"github.com/streadway/amqp"
)

type WorkerMessageQueue struct {
	BaseMQ
	queueName 		string
}

func NewWorkerMessageQueue(queueName string) IMessageQueue {
	var err error
	workerRabbitMQ := &WorkerMessageQueue{}
	workerRabbitMQ.rabbitConn, err = amqp.Dial(rabbitMQUrl)
	if err != nil {
		log.Errorf(err.Error())
		panic(err)
	}
	workerRabbitMQ.channel, err = workerRabbitMQ.rabbitConn.Channel()
	if err != nil {
		log.Errorf(err.Error())
		panic(err)
	}
	workerRabbitMQ.queueName = queueName
	return workerRabbitMQ
}

func (this *WorkerMessageQueue) Publish(routingKey string, content string) error {
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
func (this *WorkerMessageQueue) Consume(bindingKey string) (<-chan amqp.Delivery, error) {
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