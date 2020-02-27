package mq

import (
	"errors"
	"bookzone/conf"
	"github.com/streadway/amqp"
)

type ModelType int

const (
	MQ_MODEL_SIMPLE ModelType = iota
	MQ_MODEL_WORK
	MQ_MODEL_PUBSUB
	MQ_MODEL_ROUTE
	MQ_MODEL_TOPIC
)

var rabbitMQUrl string

func init() {
	urlKey := "url"
	rabbitMQUrl = conf.GlobalCfg.Section("messagequeue").Key(urlKey).String()
	if rabbitMQUrl == "" {
		panic(errors.New("invalid rabbitmq url."))
	}
}

type BaseMQ struct {
	rabbitConn 		*amqp.Connection
	channel 		*amqp.Channel
}

type IMessageQueue interface {
	Publish(routingKey string, content string) error
	Consume(bindingKey string) (<-chan amqp.Delivery, error)
}

func NewRabbitMQ(t ModelType, queueName string, exchageName string) IMessageQueue {
	switch t {
	case MQ_MODEL_SIMPLE:
		return NewSimpleMessageQueue(queueName)
	case MQ_MODEL_WORK:
		return NewWorkerMessageQueue(queueName)
	case MQ_MODEL_PUBSUB:
		return NewPubSubMessageQueue(exchageName)
	case MQ_MODEL_ROUTE:
		return NewRouteMessageQueue(exchageName)
	case MQ_MODEL_TOPIC:
		return NewTopicMessageQueue(exchageName)
	default:
		return NewSimpleMessageQueue(queueName)
	}
}