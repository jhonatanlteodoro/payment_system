package ports

import (
	"context"
	"github.com/jhonatanlteodoro/payment_system/src/types"
)

type ProcessingRules interface {
	ProcessRules(ctx context.Context, payment *types.Payment, summary *types.AccountPaymentQuarterlySummary) (string, error)
}
