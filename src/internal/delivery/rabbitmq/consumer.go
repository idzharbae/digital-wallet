package rabbitmq

import (
	"github.com/idzharbae/digital-wallet/src/internal/constants"
	"github.com/idzharbae/digital-wallet/src/internal/gateway/rmq_gateway"
	"github.com/idzharbae/digital-wallet/src/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/palantir/stacktrace"
	"github.com/redis/go-redis/v9"
)

type RabbitMQConsumer struct {
	consumerConn *rmq_gateway.RMQConsumer
	dbConn       *pgxpool.Pool
	redisClient  *redis.Client

	userTransactionRepo repository.UserTransactionRepository
}

func NewConsumer(consumerConn *rmq_gateway.RMQConsumer, dbConn *pgxpool.Pool, redisclient *redis.Client, userTransactionRepo repository.UserTransactionRepository) *RabbitMQConsumer {
	consumer := &RabbitMQConsumer{
		consumerConn:        consumerConn,
		dbConn:              dbConn,
		redisClient:         redisclient,
		userTransactionRepo: userTransactionRepo,
	}
	consumer.consumerConn.AddHandler(constants.EXAMPLE_QUEUE, consumer.HandleExample)
	consumer.consumerConn.AddHandler(constants.TOTAL_DEBIT_QUEUE, consumer.HandleTotalDebit)
	return consumer
}

func (c *RabbitMQConsumer) Consume() error {
	if err := c.consumerConn.Consume(); err != nil {
		return stacktrace.Propagate(err, "failed to start consumer")
	}

	return nil
}
