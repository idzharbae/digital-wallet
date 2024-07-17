package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/idzharbae/digital-wallet/src/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/palantir/stacktrace"
)

type User struct {
	userTokenRepository   repository.UserTokenRepository
	userBalanceRepository repository.UserBalanceRepository
	transactionHandler    repository.TransactionHandler
}

func NewUser(
	userTokenRepository repository.UserTokenRepository,
	userBalanceRepository repository.UserBalanceRepository,
	transactionHandler repository.TransactionHandler) UserUC {
	return &User{
		userTokenRepository:   userTokenRepository,
		userBalanceRepository: userBalanceRepository,
		transactionHandler:    transactionHandler,
	}
}

func (u *User) RegisterUser(ctx context.Context, username string) (string, error) {
	token := uuid.New()

	err := u.transactionHandler.ExecuteTransaction(ctx, func(tx pgx.Tx) error {
		err := u.userTokenRepository.
			WithTransaction(tx).
			InsertUserToken(ctx, username, token.String())
		if err != nil {
			return stacktrace.Propagate(err, "RegisterUser Transaction: failed to insert token")
		}

		err = u.userBalanceRepository.
			WithTransaction(tx).
			CreateUserBalance(ctx, username)
		if err != nil {
			return stacktrace.Propagate(err, "RegisterUser Transaction: failed to create user balance")
		}

		return nil
	})
	if err != nil {
		return "", stacktrace.Propagate(err, "RegisterUser: error occured during transaction process")
	}

	return token.String(), nil
}
