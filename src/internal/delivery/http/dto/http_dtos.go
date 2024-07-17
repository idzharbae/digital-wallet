package dto

type RegisterUserRequest struct {
	Username string `json:"username"`
}

type BalanceTopUpRequest struct {
	Amount int `json:"amount"`
}

type TransferRequest struct {
	ToUsername string `json:"to_username"`
	Amount     int    `json:"amount"`
}
