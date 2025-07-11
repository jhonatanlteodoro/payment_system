package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"html"
	"log"
	"net/http"
	"os"
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

func main() {
	ctx := context.Background()
	ConnectDB(ctx)
	defer dbConn.Close(ctx)

	http.HandleFunc("/api/v1/payment", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		data := &Transaction{}
		if err := json.NewDecoder(r.Body).Decode(data); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		fmt.Println("Amount ", data.Amount.String())

		balanceFromAccountQuery := `
			SELECT amount FROM balances WHERE account_id=$1 LIMIT 1;
		`

		balanceFromAccount := &Balance{}
		if err := dbConn.QueryRow(ctx, balanceFromAccountQuery, data.FromAccount).Scan(
			&balanceFromAccount.Amount,
		); err != nil {
			http.Error(w, "Unable to query database", http.StatusInternalServerError)
			return
		}

		fmt.Println("Balance: ", balanceFromAccount.Amount.String())

		balanceToAccountQuery := `
			SELECT amount FROM balances WHERE account_id=$1 LIMIT 1;
		`

		balanceToAccount := &Balance{}
		if err := dbConn.QueryRow(ctx, balanceToAccountQuery, data.ToAccount).Scan(
			&balanceToAccount.Amount,
		); err != nil {
			http.Error(w, "Unable to query database", http.StatusInternalServerError)
			return
		}

		fmt.Println("Balance: ", balanceToAccount.Amount.String())

		id, err := uuid.NewV7()
		if err != nil {
			http.Error(w, "Unable to generate UUID", http.StatusInternalServerError)
			return
		}

		query := `
			INSERT INTO transactions (id, from_account_id, to_account_id, amount, description) VALUES ($1, $2, $3, $4, $5);
		`
		if _, err := dbConn.Exec(ctx, query, id.String(), data.FromAccount, data.ToAccount, &data.Amount, data.Description); err != nil {
			log.Printf("query failed: %v\n", err)
			http.Error(w, "Unable to insert transaction", http.StatusInternalServerError)
			return
		}

		fromAccountNewBalanceQuery := `
        	UPDATE balances SET amount=$2 where account_id=$1;
        `
		newBalance := balanceFromAccount.Amount.Sub(data.Amount.Decimal)
		if _, err = dbConn.Exec(ctx, fromAccountNewBalanceQuery, data.FromAccount, &newBalance); err != nil {
			log.Printf("query failed: %v\n", err)
			http.Error(w, "Unable to insert transaction", http.StatusInternalServerError)
			return
		}

		newBalanceToAccount := balanceToAccount.Amount.Add(data.Amount.Decimal)
		if _, err = dbConn.Exec(ctx, fromAccountNewBalanceQuery, data.ToAccount, &newBalanceToAccount); err != nil {
			log.Printf("query failed: %v\n", err)
			http.Error(w, "Unable to insert transaction", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
