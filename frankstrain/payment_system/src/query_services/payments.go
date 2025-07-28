package query_services

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jhonatanlteodoro/payment_system/src/ports"
	"github.com/jhonatanlteodoro/payment_system/src/types"
)

type PaymentsQuery struct{}

func NewPaymentsQuery() ports.PaymentsQuery {
	return &PaymentsQuery{}
}

func (m *PaymentsQuery) CreatePayment(ctx context.Context, payment *types.Payment, onGoingDBTransaction pgx.Tx) error {
	query := `
    	INSERT INTO payments (id, from_account_id, to_account_id, amount, description, status) VALUES ($1, $2, $3, $4, $5, $6)
	`
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	payment.ID = id.String()

	_, err = onGoingDBTransaction.Exec(ctx, query, payment.ID, payment.FromAccount, payment.ToAccount, payment.Amount, payment.Description, payment.Status)
	return err
}

func (m *PaymentsQuery) UpdatePaymentStatus(ctx context.Context, payment *types.Payment, onGoingDBTransaction pgx.Tx) error {
	query := `UPDATE payments SET status = $1 WHERE id = $2`
	_, err := onGoingDBTransaction.Exec(ctx, query, payment.Status, payment.ID)

	return err
}
