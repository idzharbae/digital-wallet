package repository

//go:generate mockgen -destination=repomock/usertokenrepo_mock.go -package=repomock github.com/idzharbae/digital-wallet/src/internal/repository UserTokenRepository
type UserTokenRepository interface {
	InsertUserToken(username, token string) error
}
