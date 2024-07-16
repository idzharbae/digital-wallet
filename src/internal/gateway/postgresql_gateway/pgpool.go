package postgresql_gateway

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/palantir/stacktrace"
)

func NewPgPool(connectionString string) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Unable to parse DATABASE_URL")
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Unable to create connection pool")
	}

	return pool, nil
}
