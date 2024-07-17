package gateway

import (
	"github.com/idzharbae/digital-wallet/src/internal/gateway/rmq_gateway"
)

type rabbitMqGateway struct {
	producer *rmq_gateway.RMQProducer
}

func NewRabbitMqGateway(producer *rmq_gateway.RMQProducer) RabbitMqGateway {
	return &rabbitMqGateway{
		producer: producer,
	}
}

func (rmg *rabbitMqGateway) PublishMessage(queueName, contentType string, body []byte) error {
	return rmg.producer.PublishMessage(queueName, contentType, body)
}
