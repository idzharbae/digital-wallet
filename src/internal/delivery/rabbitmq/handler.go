package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/idzharbae/digital-wallet/src/internal/delivery/rabbitmq/dto"
	"github.com/palantir/stacktrace"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

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
