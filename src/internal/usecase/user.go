package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/idzharbae/digital-wallet/src/internal/repository"
	"github.com/palantir/stacktrace"
)

type User struct {
	userTokenRepository repository.UserTokenRepository
}

func NewUser(userTokenRepository repository.UserTokenRepository) UserUC {
	return &User{
		userTokenRepository: userTokenRepository,
	}
}

func (u *User) RegisterUser(ctx context.Context, username string) (string, error) {
	token := uuid.New()
	if err := u.userTokenRepository.InsertUserToken(ctx, username, token.String()); err != nil {
		return "", stacktrace.Propagate(err, "RegisterUser: failed to insert token")
	}

	return token.String(), nil
}
