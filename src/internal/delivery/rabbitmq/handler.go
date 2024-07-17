package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/idzharbae/digital-wallet/src/internal/constants"
	"github.com/idzharbae/digital-wallet/src/internal/delivery/rabbitmq/dto"
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

func (c *RabbitMQConsumer) HandleTotalDebit(queue string, msg amqp.Delivery) error {
	var request dto.DebitMessage
	if err := json.Unmarshal(msg.Body, &request); err != nil {
		return stacktrace.Propagate(err, "HandleTotalDebit: invalid message: %s", string(msg.Body))
	}

	if request.Username == "" {
		return stacktrace.NewError("Handle Total Debit: username is empty")
	}

	err := c.userTransactionRepo.UpsertTotalDebit(context.Background(), request.Username, request.Amount)
	if err != nil {
		return stacktrace.Propagate(err, "HandleTotalDebit: failed to upsert total debit for user %s", request.Username)
	}

	log.Info().Msg(fmt.Sprintf("Updated user %s total_debit by %d", request.Username, request.Amount))

	return nil
}
