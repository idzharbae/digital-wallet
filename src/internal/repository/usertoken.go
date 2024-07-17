package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/palantir/stacktrace"
	"github.com/redis/go-redis/v9"
)

type UserToken struct {
	dbConn      *pgxpool.Pool
	redisClient *redis.Client
}

func NewUserToken(dbConn *pgxpool.Pool, redisClient *redis.Client) UserTokenRepository {
	return &UserToken{
		dbConn:      dbConn,
		redisClient: redisClient,
	}
}

func (ut *UserToken) InsertUserToken(ctx context.Context, username string, token string) error {
	_, err := ut.dbConn.Exec(ctx, `INSERT INTO user_token(username, token) VALUES ($1, $2);`, username, token)
	if err != nil {
		return stacktrace.Propagate(err, "InsertUserToken: failed to insert to db for username: %s", username)
	}

	return nil
}
