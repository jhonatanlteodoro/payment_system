package query_services

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jhonatanlteodoro/payment_system/src/ports"
	"github.com/jhonatanlteodoro/payment_system/src/types"
)

type QuarterlyAccountSummary struct{}

func NewQuarterlyAccountSummary() ports.QuarterlyAccountSummary {
	return &QuarterlyAccountSummary{}
}

func (m *QuarterlyAccountSummary) GetSummary(ctx context.Context, accountId string, onGoingDBTransaction pgx.Tx) (*types.AccountPaymentQuarterlySummary, error) {
	query := `SELECT * FROM account_payments_quarterly_summary WHERE account_id = $1`
	data := &types.AccountPaymentQuarterlySummary{}
	err := onGoingDBTransaction.QueryRow(ctx, query, accountId).Scan(data)
	return data, err
}
