package rabbitmq

import (
	"fmt"
	"github.com/ahaostudy/calendar_reminder/conf"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"sync"
)

var (
	once sync.Once
	// global connection object
	conn *amqp.Connection
)

func initConn() {
	mqUrl := fmt.Sprintf("amqp://%s:%s@%s/%s",
		conf.GetConf().RabbitMQ.Username,
		conf.GetConf().RabbitMQ.Password,
		conf.GetConf().RabbitMQ.Address,
		conf.GetConf().RabbitMQ.VHost)

	var err error
	conn, err = amqp.Dial(mqUrl)
	if err != nil {
		logrus.Fatalf("failed to connect rabbitmq: %s", err.Error())
	}
}

type RabbitMQ struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	Exchange string
	Key      string
}

func NewRabbitMQ(exchange string, key string) *RabbitMQ {
	return &RabbitMQ{Exchange: exchange, Key: key}
}

func NewWorkRabbitMQ(queue string) *RabbitMQ {
	// new rabbitmq
	rabbitmq := NewRabbitMQ("", queue)

	// get connection
	once.Do(initConn)
	rabbitmq.conn = conn

	// get channel
	var err error
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	if err != nil {
		panic(fmt.Errorf("failed to open a channel: %s\n", err.Error()))
	}

	return rabbitmq
}

func (r *RabbitMQ) Destroy() {
	if r == nil {
		return
	}
	if r.conn != nil {
		_ = r.conn.Close()
	}
	if r.channel != nil {
		_ = r.channel.Close()
	}
}

func (r *RabbitMQ) Publish(message []byte) error {
	_, err := r.channel.QueueDeclare(r.Key, false, false, false, false, nil)
	if err != nil {
		return err
	}

	return r.channel.Publish(r.Exchange, r.Key, false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
}

func (r *RabbitMQ) Consume(handle func(msg *amqp.Delivery) error) {
	q, err := r.channel.QueueDeclare(r.Key, false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	msgs, err := r.channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	for msg := range msgs {
		if err := handle(&msg); err != nil {
			logrus.Error(err.Error())
		}
	}
}
