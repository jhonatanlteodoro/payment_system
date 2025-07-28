package query_services

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jhonatanlteodoro/payment_system/src/ports"
	"github.com/jhonatanlteodoro/payment_system/src/types"
)

type BalanceQuery struct{}

func NewBalanceQuery() ports.BalanceQuery {
	return &BalanceQuery{}
}

func (m *BalanceQuery) GetBalance(ctx context.Context, accountId string, onGoingDBTransaction pgx.Tx) (*types.Balance, error) {
	data := &types.Balance{}
	query := `SELECT * FROM balances WHERE account_id = $1`
	err := onGoingDBTransaction.QueryRow(ctx, query, accountId).Scan(data)
	return data, err
}

func (m *BalanceQuery) UpdateBalance(ctx context.Context, balance *types.Balance, onGoingDBTransaction pgx.Tx) error {
	query := `UPDATE accounts SET amount = $1 WHERE id = $2`
	_, err := onGoingDBTransaction.Exec(ctx, query, balance.Amount, balance.AccountID)
	return err
}
