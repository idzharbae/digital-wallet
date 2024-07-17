package usecase

import (
	"context"

	"github.com/idzharbae/digital-wallet/src/internal/entity"
)

//go:generate mockgen -destination=ucmock/useruc_mock.go -package=ucmock github.com/idzharbae/digital-wallet/src/internal/usecase UserUC
type UserUC interface {
	RegisterUser(ctx context.Context, username string) (string, error)
	GetUserNameFromToken(ctx context.Context, token string) (string, error)
	GetUserBalance(ctx context.Context, username string) (int, error)
	TopUpUserBalance(ctx context.Context, username string, topUpAmount int) (int, error)
}

type TransactionUC interface {
	TransferBalance(ctx context.Context, senderUsername, recipientUsername string, trasnferAmount int) error
	GetTopTransactingUsers(ctx context.Context) ([]entity.TotalDebit, error)
	GetUserTopTransactions(ctx context.Context, username string) ([]entity.UserTransaction, error)
}
