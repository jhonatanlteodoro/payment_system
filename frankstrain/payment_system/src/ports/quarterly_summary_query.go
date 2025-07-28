package ports

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jhonatanlteodoro/payment_system/src/types"
)

type QuarterlyAccountSummary interface {
	GetSummary(ctx context.Context, accountId string, onGoingDBTransaction pgx.Tx) (*types.AccountPaymentQuarterlySummary, error)
}
