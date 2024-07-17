package repository

import (
	"context"

	"github.com/idzharbae/digital-wallet/src/internal/entity"
	"github.com/jackc/pgx/v5"
)

//go:generate mockgen -destination=repomock/pgx_tx_mock.go -package=repomock github.com/jackc/pgx/v5 Tx
type TransactionFunction func(trx pgx.Tx) error

//go:generate mockgen -destination=repomock/transactionhandler_mock.go -package=repomock github.com/idzharbae/digital-wallet/src/internal/repository TransactionHandler
type TransactionHandler interface {
	ExecuteTransaction(ctx context.Context, f TransactionFunction) error
}

//go:generate mockgen -destination=repomock/usertokenrepo_mock.go -package=repomock github.com/idzharbae/digital-wallet/src/internal/repository UserTokenRepository
type UserTokenRepository interface {
	InsertUserToken(ctx context.Context, username, token string) error
	GetUserNameByToken(ctx context.Context, token string) (string, error)

	WithTransaction(tx pgx.Tx) UserTokenRepository
}

//go:generate mockgen -destination=repomock/userbalancerepo_mock.go -package=repomock github.com/idzharbae/digital-wallet/src/internal/repository UserBalanceRepository
type UserBalanceRepository interface {
	CreateUserBalance(ctx context.Context, username string) error
	GetUserBalance(ctx context.Context, username string) (int, error)
	UpdateBalance(ctx context.Context, username string, balance int) error
	GetUserBalanceForUpdate(ctx context.Context, username string) (int, error)

	WithTransaction(tx pgx.Tx) UserBalanceRepository
}

type UserTransactionRepository interface {
	InsertTransaction(ctx context.Context, username, secondPartyUsername string, transactionType entity.TransactionType, amount int) error
	UpsertTotalDebit(ctx context.Context, username string, amount int) error
	GetTopTransactingUsers(ctx context.Context) ([]entity.TotalDebit, error)
	RefreshTopTransactingUsers(ctx context.Context) ([]entity.TotalDebit, error)
	GetUserTopTransactions(ctx context.Context, username string) ([]entity.UserTransaction, error)

	WithTransaction(tx pgx.Tx) UserTransactionRepository
}
