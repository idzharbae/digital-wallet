package repository

import (
	"context"

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
	WithTransaction(tx pgx.Tx) UserTokenRepository
}

//go:generate mockgen -destination=repomock/userbalancerepo_mock.go -package=repomock github.com/idzharbae/digital-wallet/src/internal/repository UserBalanceRepository
type UserBalanceRepository interface {
	CreateUserBalance(ctx context.Context, username string) error
	WithTransaction(tx pgx.Tx) UserBalanceRepository
}
