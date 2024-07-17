package dto

type RegisterUserRequest struct {
	Username string `json:"username"`
}

type BalanceTopUpRequest struct {
	Amount int `json:"amount"`
}
