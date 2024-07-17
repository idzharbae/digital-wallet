package gateway

type RabbitMqGateway interface {
	PublishMessage(queueName, contentType string, body []byte) error
}
