package query_services

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jhonatanlteodoro/payment_system/src/ports"
	"github.com/jhonatanlteodoro/payment_system/src/types"
)

type AccountQuery struct{}

func NewAccountQuery() ports.AccountQuery {
	return &AccountQuery{}
}

func (m *AccountQuery) GetAccount(ctx context.Context, id string, onGoingDBTransaction pgx.Tx) (*types.Account, error) {
	acc := &types.Account{}
	err := onGoingDBTransaction.QueryRow(ctx, "SELECT * FROM accounts WHERE id = $1", id).Scan(acc)

	return acc, err
}
