package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/palantir/stacktrace"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type userBalance struct {
	dbConn      Pgx
	redisClient *redis.Client
}

const (
	UserBalanceKey           = "user_balance:%s"
	UserBalanceCacheDuration = time.Minute
)

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

func (ub *userBalance) GetUserBalance(ctx context.Context, username string) (int, error) {
	var balance int
	redisKey := fmt.Sprintf(UserBalanceKey, username)
	val, err := ub.redisClient.Get(ctx, redisKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &balance)
		return balance, nil
	}
	if err != redis.Nil {
		log.Error().Err(err).Msg(fmt.Sprintf("GetUserBalance: failed to get balance from redis for username %s", username))
	}

	err = ub.dbConn.QueryRow(ctx, `SELECT balance FROM user_balance WHERE username = $1`, username).Scan(&balance)
	if err != nil {
		return 0, stacktrace.Propagate(err, "GetUserBalance: failed to read user balance for username %s", username)
	}

	err = ub.redisClient.Set(ctx, redisKey, balance, UserBalanceCacheDuration).Err()
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("GetUserBalance: failed to set redis value for username %s", username))
	}

	return balance, nil
}

func (ub *userBalance) GetUserBalanceForUpdate(ctx context.Context, username string) (int, error) {
	var balance int
	err := ub.dbConn.QueryRow(ctx, `SELECT balance FROM user_balance WHERE username = $1 FOR UPDATE`, username).Scan(&balance)
	if err != nil {
		return 0, stacktrace.Propagate(err, "GetUserBalanceForUpdate: failed to read user balance for username %s", username)
	}

	return balance, nil
}

func (ub *userBalance) UpdateBalance(ctx context.Context, username string, balance int) error {
	_, err := ub.dbConn.Exec(ctx, `UPDATE user_balance SET balance = $1 WHERE username = $2`, balance, username)
	if err != nil {
		return stacktrace.Propagate(err, "UpdateBalance: failed to update user balance for username %s", username)
	}

	// Clear redis
	redisKey := fmt.Sprintf(UserBalanceKey, username)
	err = ub.redisClient.Del(ctx, redisKey).Err()
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("UpdateBalance: failed to clear user balance redis for user %s", username))
	}

	return nil
}

func (ub *userBalance) WithTransaction(tx pgx.Tx) UserBalanceRepository {
	return &userBalance{
		dbConn:      tx,
		redisClient: ub.redisClient,
	}
}
