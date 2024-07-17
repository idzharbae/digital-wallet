package usecase

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/idzharbae/digital-wallet/src/internal/constants"
	"github.com/idzharbae/digital-wallet/src/internal/entity"
	"github.com/idzharbae/digital-wallet/src/internal/gateway"
	"github.com/idzharbae/digital-wallet/src/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/palantir/stacktrace"
)

var (
	ErrNotEnoughBalance      = errors.New("sender's balance is lower than transfer amount")
	ErrInvalidTransferAmount = errors.New("transfer amount must be positive integer")
)

type transaction struct {
	transactionHandler  repository.TransactionHandler
	userTransactionRepo repository.UserTransactionRepository
	userBalanceRepo     repository.UserBalanceRepository
	rmqPublisher        gateway.RabbitMqGateway
}

func NewTransaction(
	transactionHandler repository.TransactionHandler,
	userTransactionRepo repository.UserTransactionRepository,
	userBalanceRepo repository.UserBalanceRepository,
	rmqPublisher gateway.RabbitMqGateway) TransactionUC {
	return &transaction{
		userTransactionRepo: userTransactionRepo,
		rmqPublisher:        rmqPublisher,
		userBalanceRepo:     userBalanceRepo,
		transactionHandler:  transactionHandler,
	}
}

func (t *transaction) TransferBalance(ctx context.Context, senderUsername, recipientUsername string, transferAmount int) error {
	if transferAmount <= 0 {
		return ErrInvalidTransferAmount
	}
	err := t.transactionHandler.ExecuteTransaction(ctx, func(trx pgx.Tx) error {
		userBalanceTrx := t.userBalanceRepo.WithTransaction(trx)
		deductSenderBalance := func() error {
			senderBalance, err := userBalanceTrx.GetUserBalanceForUpdate(ctx, senderUsername)
			if err != nil {
				return stacktrace.Propagate(err, "TransferBalance Transaction: failed to get sender's balance: %s", senderUsername)
			}

			if senderBalance < transferAmount {
				return ErrNotEnoughBalance
			}

			err = userBalanceTrx.UpdateBalance(ctx, senderUsername, senderBalance-transferAmount)
			if err != nil {
				return stacktrace.Propagate(err, "TransferBalance Transaction: failed to deduct sender's balance: %s", senderUsername)
			}

			return nil
		}

		increaseRecipientBalance := func() error {
			recipientBalance, err := userBalanceTrx.GetUserBalanceForUpdate(ctx, recipientUsername)
			if err != nil {
				return stacktrace.Propagate(err, "TransferBalance Transaction: failed to get recipient's balance: %s", recipientUsername)
			}

			err = userBalanceTrx.UpdateBalance(ctx, recipientUsername, recipientBalance+transferAmount)
			if err != nil {
				return stacktrace.Propagate(err, "TransferBalance Transaction: failed to add recipient's balance: %s", recipientUsername)
			}

			return nil
		}

		// Avoid deadlock by executing balance update by username order
		if senderUsername < recipientUsername {
			err := deductSenderBalance()
			if err != nil {
				return err
			}
			err = increaseRecipientBalance()
			if err != nil {
				return err
			}
		} else {
			err := increaseRecipientBalance()
			if err != nil {
				return err
			}
			err = deductSenderBalance()
			if err != nil {
				return err
			}
		}

		// Insert transaction record
		userTransactionTrx := t.userTransactionRepo.WithTransaction(trx)
		err := userTransactionTrx.InsertTransaction(ctx, senderUsername, recipientUsername, entity.TransactionTypeDebit, transferAmount)
		if err != nil {
			return stacktrace.Propagate(err, "TransferBalance Transaction: failed to insert debit transaction for user: %s", senderUsername)
		}

		err = userTransactionTrx.InsertTransaction(ctx, recipientUsername, senderUsername, entity.TransactionTypeCredit, transferAmount)
		if err != nil {
			return stacktrace.Propagate(err, "TransferBalance Transaction: failed to insert credit transaction for user: %s", senderUsername)
		}

		rmqPayload, err := json.Marshal(map[string]any{
			"username": senderUsername,
			"amount":   transferAmount,
		})
		if err != nil {
			return stacktrace.Propagate(err, "TransferBalance Transaction: failed to marshal debit transaction info for user: %s", senderUsername)
		}
		err = t.rmqPublisher.PublishMessage(constants.TOTAL_DEBIT_QUEUE, "application/json", rmqPayload)
		if err != nil {
			return stacktrace.Propagate(err, "TransferBalance Transaction: failed to marshal debit transaction info for user: %s", senderUsername)
		}

		return nil
	})
	if err != nil {
		return stacktrace.Propagate(err, "TransferBalance: error occured during transaction")
	}

	return nil
}

func (t *transaction) GetTopTransactingUsers(ctx context.Context) ([]entity.TotalDebit, error) {
	return t.userTransactionRepo.GetTopTransactingUsers(ctx)
}

func (t *transaction) GetUserTopTransactions(ctx context.Context, username string) ([]entity.UserTransaction, error) {
	return t.userTransactionRepo.GetUserTopTransactions(ctx, username)
}
