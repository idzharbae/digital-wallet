package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/palantir/stacktrace"
	"github.com/redis/go-redis/v9"
)

type userBalance struct {
	dbConn      Pgx
	redisClient *redis.Client
}

func NewUserBalance(dbConn Pgx, redisClient *redis.Client) UserBalanceRepository {
	return &userBalance{
		dbConn:      dbConn,
		redisClient: redisClient,
	}
}

func (ub *userBalance) CreateUserBalance(ctx context.Context, username string) error {
	_, err := ub.dbConn.Exec(ctx, `INSERT INTO user_balance(username, balance) VALUES ($1, $2);`, username, 0)
	if err != nil {
		return stacktrace.Propagate(err, "CreateUserBalance: failed to insert to db for username: %s", username)
	}

	return nil
}

func (ub *userBalance) WithTransaction(tx pgx.Tx) UserBalanceRepository {
	return &userBalance{
		dbConn:      tx,
		redisClient: ub.redisClient,
	}
}
