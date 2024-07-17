package rmq_gateway

import (
	"github.com/palantir/stacktrace"

	"github.com/streadway/amqp"
)

type RMQProducer struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewProducer(connectionString string) (*RMQProducer, error) {
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to connect to rabbitMQ")
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to open rabbitMQ channel")
	}

	return &RMQProducer{
		conn: conn,
		ch:   ch,
	}, nil
}

func (x *RMQProducer) PublishMessage(queueName string, contentType string, body []byte) error {
	q, err := x.ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return stacktrace.Propagate(err, "PublishMessage: failed to declare queue")
	}

	err = x.ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: contentType,
			Body:        body,
		})
	if err != nil {
		return stacktrace.Propagate(err, "PublishMessage: failed to publish message")
	}

	return nil
}
