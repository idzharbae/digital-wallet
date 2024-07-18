package gateway

//go:generate mockgen -destination=gatewaymock/rabbitmqgateway_mock.go -package=gatewaymock github.com/idzharbae/digital-wallet/src/internal/gateway RabbitMqGateway
type RabbitMqGateway interface {
	PublishMessage(queueName, contentType string, body []byte) error
}
