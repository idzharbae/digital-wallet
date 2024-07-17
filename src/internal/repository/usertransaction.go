package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/idzharbae/digital-wallet/src/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/palantir/stacktrace"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

const (
	TotalDebitKey            = "totaldebit"
	UserTopTransactionKey    = "user_top_trx:%s"
	UserTopTransactionExpiry = time.Minute
)

type userTransaction struct {
	dbConn      Pgx
	redisClient *redis.Client
}

func NewUserTransaction(dbConn Pgx, redisClient *redis.Client) UserTransactionRepository {
	return &userTransaction{
		dbConn:      dbConn,
		redisClient: redisClient,
	}
}

func (ut *userTransaction) InsertTransaction(
	ctx context.Context, username, secondPartyUsername string, transactionType entity.TransactionType, amount int,
) error {
	_, err := ut.dbConn.Exec(
		ctx, `INSERT INTO transactions (username, second_party, amount, "type") VALUES ($1, $2, $3, $4)`, username, secondPartyUsername, amount, transactionType,
	)
	if err != nil {
		return stacktrace.Propagate(err, "InsertTransaction: failed to insert transaction for user %s", username)
	}

	redisKey := fmt.Sprintf(UserTopTransactionKey, username)
	err = ut.redisClient.Del(ctx, redisKey).Err()
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("InsertTransaction: failed to clear user top transactions redis for user %s", username))
	}

	return nil
}

func (ut *userTransaction) UpsertTotalDebit(
	ctx context.Context, username string, amount int,
) error {
	_, err := ut.dbConn.Exec(
		ctx,
		`INSERT INTO user_total_debit (username, total_debit) VALUES ($1, $2)
			ON CONFLICT (username) DO
			UPDATE SET total_debit = user_total_debit.total_debit + excluded.total_debit;
		`, username, amount,
	)
	if err != nil {
		return stacktrace.Propagate(err, "UpsertTotalDebit: failed to upsert total debit for user %s", username)
	}

	return nil
}

func (ut *userTransaction) GetTopTransactingUsers(ctx context.Context) ([]entity.TotalDebit, error) {
	result := []entity.TotalDebit{}
	val, err := ut.redisClient.Get(ctx, TotalDebitKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &result)
		return result, nil
	}
	if err != redis.Nil {
		log.Error().Err(err).Msg("GetTopTransactingUsers: failed to get top transacting users from redis")
	}

	// Get from DB
	return ut.RefreshTopTransactingUsers(ctx)
}

func (ut *userTransaction) RefreshTopTransactingUsers(ctx context.Context) ([]entity.TotalDebit, error) {
	rows, err := ut.dbConn.Query(ctx, `SELECT username, total_debit FROM user_total_debit ORDER BY total_debit DESC LIMIT 10`)
	if err != nil {
		return nil, stacktrace.Propagate(err, "RefreshTopTransactingUsers: failed to get total debit")
	}
	defer rows.Close()

	result := []entity.TotalDebit{}
	for rows.Next() {
		var totalDebit entity.TotalDebit
		err := rows.Scan(&totalDebit.UserName, &totalDebit.Amount)
		if err != nil {
			return nil, stacktrace.Propagate(err, "RefreshTopTransactingUsers: failed to scan total debit")
		}

		result = append(result, totalDebit)
	}

	marshaledResult, err := json.Marshal(result)
	if err != nil {
		return nil, stacktrace.Propagate(err, "RefreshTopTransactingUsers: failed to marshal total debit")
	}
	err = ut.redisClient.Set(ctx, TotalDebitKey, marshaledResult, 0).Err()
	if err != nil {
		return nil, stacktrace.Propagate(err, "RefreshTopTransactingUsers: failed to refresh total debit redis")
	}

	return result, nil
}

func (ut *userTransaction) GetUserTopTransactions(ctx context.Context, username string) ([]entity.UserTransaction, error) {
	result := []entity.UserTransaction{}
	redisKey := fmt.Sprintf(UserTopTransactionKey, username)
	val, err := ut.redisClient.Get(ctx, redisKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &result)
		return result, nil
	}
	if err != redis.Nil {
		log.Error().Err(err).Msg(fmt.Sprintf("GetUserTopTransactions: failed to get top transactions from redis for user %s", username))
	}

	rows, err := ut.dbConn.Query(ctx, `SELECT second_party, amount, "type" FROM transactions WHERE username = $1 ORDER BY amount DESC LIMIT 10;`, username)
	if err != nil {
		return nil, stacktrace.Propagate(err, "GetUserTopTransactions: failed to get user top transactions: %s", username)
	}
	defer rows.Close()

	for rows.Next() {
		var transaction entity.UserTransaction
		err := rows.Scan(&transaction.UserName, &transaction.Amount, &transaction.Type)
		if err != nil {
			return nil, stacktrace.Propagate(err, "GetUserTopTransactions: failed to scan user top transactions: %s", username)
		}

		result = append(result, transaction)
	}

	marshaledResult, err := json.Marshal(result)
	if err != nil {
		return nil, stacktrace.Propagate(err, "GetUserTopTransactions: failed to marshal user top transactions: %s", username)
	}
	err = ut.redisClient.Set(ctx, redisKey, marshaledResult, UserTopTransactionExpiry).Err()
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("Failed to set top transaction redis for user %s", username))
	}

	return result, nil
}

func (ut *userTransaction) WithTransaction(tx pgx.Tx) UserTransactionRepository {
	return &userTransaction{
		dbConn:      tx,
		redisClient: ut.redisClient,
	}
}
