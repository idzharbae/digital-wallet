package dto

type DebitMessage struct {
	Username string `json:"username"`
	Amount   int    `json:"amount"`
}
