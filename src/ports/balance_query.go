package ports

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jhonatanlteodoro/payment_system/src/types"
)

type BalanceQuery interface {
	GetBalance(ctx context.Context, accountId string, onGoingDBTransaction pgx.Tx) (*types.Balance, error)
	UpdateBalance(ctx context.Context, balance *types.Balance, onGoingDBTransaction pgx.Tx) error
}
