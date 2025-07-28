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

type AccountPaymentQuarterlySummary struct {
	AccountID               string `json:"account_id"`
	TotalTransactions       int    `json:"total_transactions"`
	AvgTransactionValue     int    `json:"avg_transaction_value"`
	MaxEverTransactionValue int    `json:"max_ever_transaction_value"`

	//AvgWeeklyTransactions  int
	AvgDailyTransactions int `json:"avg_daily_transactions"`
	//AvgMonthlyTransactions int

	TodayTransactions int `json:"today_transactions"`
	//WeeklyTransactions  int
	//MonthlyTransactions int

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
