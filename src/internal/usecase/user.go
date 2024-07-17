package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/idzharbae/digital-wallet/src/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/palantir/stacktrace"
)

var (
	ErrTopUpTooLarge = errors.New("top up amount is too large")
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

func (u *User) GetUserNameFromToken(ctx context.Context, token string) (string, error) {
	return u.userTokenRepository.GetUserNameByToken(ctx, token)
}

func (u *User) GetUserBalance(ctx context.Context, username string) (int, error) {
	return u.userBalanceRepository.GetUserBalance(ctx, username)
}

func (u *User) TopUpUserBalance(ctx context.Context, username string, topUpAmount int) (int, error) {
	if topUpAmount >= 10000000 {
		return 0, ErrTopUpTooLarge
	}

	var newBalance int
	err := u.transactionHandler.ExecuteTransaction(ctx, func(tx pgx.Tx) error {
		trxRepo := u.userBalanceRepository.WithTransaction(tx)
		currentBalance, err := trxRepo.GetUserBalanceForUpdate(ctx, username)
		if err != nil {
			return stacktrace.Propagate(err, "TopUpUserBalance Transaction: failed to get user balance")
		}

		newBalance = currentBalance + topUpAmount
		err = trxRepo.UpdateBalance(ctx, username, newBalance)
		if err != nil {
			return stacktrace.Propagate(err, "TopUpUserBalance Transaction: failed to update user balance")
		}

		return nil
	})
	if err != nil {
		return 0, stacktrace.Propagate(err, "TopUpUserBalance: error occured during topup transaction")
	}

	return newBalance, nil
}
