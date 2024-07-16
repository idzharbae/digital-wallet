package rabbitmq

import (
	"context"
	"os"

	"github.com/google/uuid"
	"github.com/idzharbae/digital-wallet/src/internal/constants"
	"github.com/palantir/stacktrace"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

func (c *RabbitMQConsumer) HandleExample(queue string, msg amqp.Delivery) error {
	requestId := uuid.New().String()
	ctx := context.WithValue(context.Background(), "request_id", requestId)
	log.Info().Ctx(ctx).Msgf("Message received on '%s' queue: %s", queue, string(msg.Body))

	if string(msg.Body) == "ERROR" {
		return stacktrace.NewError("error!!!")
	}

	if _, err := c.dbConn.Exec(ctx, `INSERT INTO dummy_table(message) VALUES($1)`, string(msg.Body)); err != nil {
		log.Error().Ctx(ctx).Err(err).Msg("Unable to connect to database")
		os.Exit(1)
	}

	c.redisClient.Del(ctx, constants.EXAMPLE_REDIS_KEY)

	return nil
}
