package usecase

//go:generate mockgen -destination=ucmock/useruc_mock.go -package=ucmock github.com/idzharbae/digital-wallet/src/internal/usecase UserUC
type UserUC interface {
	RegisterUser(username string) (string, error)
}
