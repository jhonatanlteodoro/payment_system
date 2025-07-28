package ports

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jhonatanlteodoro/payment_system/src/types"
)

type AccountQuery interface {
	GetAccount(ctx context.Context, id string, onGoingDBTransaction pgx.Tx) (*types.Account, error)
}
