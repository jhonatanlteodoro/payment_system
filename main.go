package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
	"html"
	"log"
	"net/http"
	"os"
	"time"
)

var dbConn *pgx.Conn

func ConnectDB(ctx context.Context) {
	dbUrl := "postgres://secret_user:secret_password@localhost:5432/payment"
	var err error
	dbConn, err = pgx.Connect(ctx, dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	if errPing := dbConn.Ping(ctx); errPing != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", errPing)
		os.Exit(1)
	}

	fmt.Println("Db connected")

	_, err = dbConn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS accounts (
			id varchar(36) PRIMARY KEY,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE TABLE IF NOT EXISTS transactions (
			id varchar(36) PRIMARY KEY,
			from_account_id varchar(36) NOT NULL,
			to_account_id varchar(36) NOT NULL,
			amount decimal(10,2) NOT NULL,
			description varchar(255) NOT NULL,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		
			FOREIGN KEY (from_account_id) REFERENCES accounts(id),
			FOREIGN KEY (to_account_id) REFERENCES accounts(id)
		);
		
		CREATE TABLE IF NOT EXISTS balances (
			id varchar(36) PRIMARY KEY,
			account_id varchar(36) NOT NULL,
			amount decimal(10,2) NOT NULL,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		
			FOREIGN KEY (account_id) REFERENCES accounts(id)
		);
	`)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create payments table: %v\n", err)
	}
}

type Transaction struct {
	ID          string
	FromAccount string
	ToAccount   string
	Amount      decimal.Decimal
	Description string
	CreatedAt   time.Time
}

type Account struct {
	ID        string
	CreatedAt time.Time
}

type Balance struct {
	ID        string
	AccountID string
	Amount    decimal.Decimal
	CreatedAt time.Time
	UpdatedAt time.Time
}

func main() {
	ctx := context.Background()
	ConnectDB(ctx)
	defer dbConn.Close(ctx)

	http.HandleFunc("/api/v1/payment", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
