package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/palantir/stacktrace"
	"github.com/rs/zerolog/log"
)

type transactionHandler struct {
	dbConn *pgxpool.Pool
}

func NewTransactionHandler(dbConn *pgxpool.Pool) TransactionHandler {
	return &transactionHandler{
		dbConn: dbConn,
	}
}

func (th *transactionHandler) ExecuteTransaction(ctx context.Context, f TransactionFunction) error {
	tx, err := th.dbConn.Begin(ctx)
	if err != nil {
		return stacktrace.Propagate(err, "ExecuteTransaction: failed to begin transaction")
	}
	defer func() {
		errRollback := tx.Rollback(ctx)
		if errRollback != nil && errRollback != pgx.ErrTxClosed {
			log.Error().Err(err).Msg("ExecuteTransaction: failed to rollback transaction")
		}
	}()

	err = f(tx)
	if err != nil {
		return stacktrace.Propagate(err, "ExecuteTransaction: error from transaction process")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return stacktrace.Propagate(err, "ExecuteTransaction: failed to commit transaction")
	}
	return nil
}
