package types

import "time"

type Payment struct {
	ID          string
	FromAccount string `json:"from_account"`
	ToAccount   string `json:"to_account"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Account struct {
	ID        string
	CreatedAt time.Time
}

type Balance struct {
	ID        string
	AccountID string
	Amount    int
	CreatedAt time.Time
	UpdatedAt time.Time
}
