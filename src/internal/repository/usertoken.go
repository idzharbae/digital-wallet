package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/palantir/stacktrace"
	"github.com/redis/go-redis/v9"
)

type userToken struct {
	dbConn      Pgx
	redisClient *redis.Client
}

func NewUserToken(dbConn Pgx, redisClient *redis.Client) UserTokenRepository {
	return &userToken{
		dbConn:      dbConn,
		redisClient: redisClient,
	}
}

func (ut *userToken) InsertUserToken(ctx context.Context, username string, token string) error {
	_, err := ut.dbConn.Exec(ctx, `INSERT INTO user_token(username, token) VALUES ($1, $2);`, username, token)
	if err != nil {
		return stacktrace.Propagate(err, "InsertUserToken: failed to insert to db for username: %s", username)
	}

	return nil
}

func (ut *userToken) GetUserNameByToken(ctx context.Context, token string) (string, error) {
	var username string
	err := ut.dbConn.QueryRow(ctx, `SELECT username FROM user_token WHERE token = $1`, token).Scan(&username)
	if err != nil {
		return "", stacktrace.Propagate(err, "GetUserNameByToken: failed to read username from db")
	}

	return username, nil
}

func (ub *userToken) WithTransaction(tx pgx.Tx) UserTokenRepository {
	return &userToken{
		dbConn:      tx,
		redisClient: ub.redisClient,
	}
}
