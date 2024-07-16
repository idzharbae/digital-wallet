package rmq_gateway

import (
	"github.com/palantir/stacktrace"
	"github.com/rs/zerolog/log"

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

func (x *RMQProducer) OnError(queueName string, err error, msg string) {
	if err != nil {
		log.Err(err).Msgf("Error occurred while publishing message on '%s' queue. Error message: %s", queueName, msg)
	}
}

func (x *RMQProducer) PublishMessage(queueName string, contentType string, body []byte) {
	q, err := x.ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	x.OnError(queueName, err, "Failed to declare a queue")

	err = x.ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: contentType,
			Body:        body,
		})
	x.OnError(queueName, err, "Failed to publish a message")
}
