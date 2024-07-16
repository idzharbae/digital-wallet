package rmq_gateway

import (
	"fmt"

	"github.com/palantir/stacktrace"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

type rabbitMqHandler func(queue string, msg amqp.Delivery) error

type RMQConsumer struct {
	conn        *amqp.Connection
	ch          *amqp.Channel
	msgHandlers map[string]rabbitMqHandler
}

func NewConsumer(connectionString string) (*RMQConsumer, error) {
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to dial RabbitMQ")
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to open RabbitMQ channel")
	}

	return &RMQConsumer{
		ch:          ch,
		conn:        conn,
		msgHandlers: make(map[string]rabbitMqHandler),
	}, nil
}

func (x *RMQConsumer) AddHandler(queueName string, handler rabbitMqHandler) {
	x.msgHandlers[queueName] = handler
}

func (x *RMQConsumer) Consume() error {
	errChan := make(chan error)
	for queueName := range x.msgHandlers {
		q, err := x.ch.QueueDeclare(
			queueName, // name
			false,     // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		if err != nil {
			return stacktrace.Propagate(err, "failed to declare queue")
		}

		msgs, err := x.ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		if err != nil {
			return stacktrace.Propagate(err, "failed to register a consumer")
		}

		go func(queueName string) {
			for d := range msgs {
				err := x.msgHandlers[queueName](queueName, d)
				if err != nil {
					errChan <- stacktrace.Propagate(err, "queue %s raised error", queueName)
				}
			}
		}(queueName)
	}

	fmt.Printf("Started listening for messages")
	for err := range errChan {
		log.Error().Err(err).Msg("RabbitMQ consumer raised error")
	}

	return nil
}
