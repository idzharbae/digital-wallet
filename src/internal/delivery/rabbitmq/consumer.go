package rabbitmq

import (
	"github.com/idzharbae/digital-wallet/src/internal/constants"
	"github.com/idzharbae/digital-wallet/src/internal/gateway/rmq_gateway"
	"github.com/idzharbae/digital-wallet/src/internal/repository"
	"github.com/palantir/stacktrace"
)

type RabbitMQConsumer struct {
	consumerConn        *rmq_gateway.RMQConsumer
	userTransactionRepo repository.UserTransactionRepository
}

func NewConsumer(consumerConn *rmq_gateway.RMQConsumer, userTransactionRepo repository.UserTransactionRepository) *RabbitMQConsumer {
	consumer := &RabbitMQConsumer{
		consumerConn:        consumerConn,
		userTransactionRepo: userTransactionRepo,
	}
	consumer.consumerConn.AddHandler(constants.TOTAL_DEBIT_QUEUE, consumer.HandleTotalDebit)
	return consumer
}

func (c *RabbitMQConsumer) Consume() error {
	if err := c.consumerConn.Consume(); err != nil {
		return stacktrace.Propagate(err, "failed to start consumer")
	}

	return nil
}
