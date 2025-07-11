package main

import (
	"database/sql/driver"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
	"time"
)

type DecimalWrapper struct {
	decimal.Decimal
}

// Scan implements the pgx Scanner interface
func (d *DecimalWrapper) Scan(src interface{}) error {
	if src == nil {
		d.Decimal = decimal.Zero
		return nil
	}

	var pgNum pgtype.Numeric
	if err := pgNum.Scan(src); err != nil {
		return fmt.Errorf("failed to scan into pgtype.Numeric: %w", err)
	}

	parsed := decimal.NewFromBigInt(pgNum.Int, pgNum.Exp)

	d.Decimal = parsed
	return nil
}

// Value implements the driver.Valuer interface for writing to database
func (d DecimalWrapper) Value() (driver.Value, error) {
	return d.String(), nil
}

type Transaction struct {
	ID          string
	FromAccount string         `json:"from_account"`
	ToAccount   string         `json:"to_account"`
	Amount      DecimalWrapper `json:"amount"`
	Description string         `json:"description"`
	CreatedAt   time.Time
}

type Account struct {
	ID        string
	CreatedAt time.Time
}

type Balance struct {
	ID        string
	AccountID string
	Amount    DecimalWrapper
	CreatedAt time.Time
	UpdatedAt time.Time
}
