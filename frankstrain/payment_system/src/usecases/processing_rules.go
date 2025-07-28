package usecases

import (
	"context"
	"github.com/jhonatanlteodoro/payment_system/src/ports"
	"github.com/jhonatanlteodoro/payment_system/src/types"
)

type ProcessingRules struct {
	minValueToTriggerVerification int
}

func NewProcessingRules(minValueToTriggerVerification int) ports.ProcessingRules {
	return &ProcessingRules{
		minValueToTriggerVerification: minValueToTriggerVerification,
	}
}

func (m *ProcessingRules) ProcessRules(ctx context.Context, payment *types.Payment, summary *types.AccountPaymentQuarterlySummary) (string, error) {

	// we will trigger the rules only if transaction is greater than X value
	if payment.Amount <= m.minValueToTriggerVerification {
		return "", nil
	}

	if summary.TotalTransactions < 10 {
		// only considering if we have data to work
		return "", nil
	}

	type checks func(p *types.Payment, s *types.AccountPaymentQuarterlySummary) (bool, string)
	for _, check := range []checks{m.isPaymentRequestSuspicious, m.isAmountSuspicious} {
		if suspicious, reason := check(payment, summary); suspicious {
			return reason, nil
		}
	}
	return "", nil
}

func (m *ProcessingRules) isAmountSuspicious(payment *types.Payment, summary *types.AccountPaymentQuarterlySummary) (bool, string) {
	if summary.MaxEverTransactionValue*5 < payment.Amount {
		return true, "Amount transaction is too high"
	}

	if summary.AvgTransactionValue*2 < payment.Amount {
		return true, "Transaction has an unusual amount"
	}

	return false, ""
}

func (m *ProcessingRules) isPaymentRequestSuspicious(_ *types.Payment, summary *types.AccountPaymentQuarterlySummary) (bool, string) {
	if summary.AvgDailyTransactions*2 < summary.TodayTransactions {
		return true, "Amount of daily transactions is twice as avg transaction"
	}

	return false, ""
}
