package ports

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jhonatanlteodoro/payment_system/src/types"
)

type PaymentsQuery interface {
	CreatePayment(ctx context.Context, payment *types.Payment, onGoingDBTransaction pgx.Tx) error
	UpdatePaymentStatus(ctx context.Context, payment *types.Payment, onGoingDBTransaction pgx.Tx) error
}
