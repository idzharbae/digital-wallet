package entity

type UserTransaction struct {
	UserName string          `json:"username"`
	Amount   int             `json:"amount"`
	Type     TransactionType `json:"type"`
}
