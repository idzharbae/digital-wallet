package repository

import "context"

//go:generate mockgen -destination=repomock/usertokenrepo_mock.go -package=repomock github.com/idzharbae/digital-wallet/src/internal/repository UserTokenRepository
type UserTokenRepository interface {
	InsertUserToken(ctx context.Context, username, token string) error
}
