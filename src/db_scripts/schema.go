package db_scripts

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
)

func CreateSchema(db *pgx.Conn) {

	// If you are planning to copy this code to have a starting pointing for a new system
	// don't do it unless if you have time to fine tuning all decisions that I took here just to make it sample be delivery it faster [for example the amount, that will not work for all cases],
	// folders design, types and etc... I understand 99% of people looking at this code wont even think on it but just in case we have some vibe coder with extra time around....
	query := `CREATE TABLE IF NOT EXISTS accounts (
				id varchar(36) PRIMARY KEY,
				created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
			);
	
			CREATE TABLE IF NOT EXISTS payments (
				id varchar(36) PRIMARY KEY,
				from_account_id varchar(36) NOT NULL,
				to_account_id varchar(36) NOT NULL,
				amount bigint NOT NULL,
				description varchar(255) NOT NULL,
				status varchar(15) NOT NULL,
				
				created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
	
				FOREIGN KEY (from_account_id) REFERENCES accounts(id),
				FOREIGN KEY (to_account_id) REFERENCES accounts(id)
			);
	
			CREATE TABLE IF NOT EXISTS balances (
				id varchar(36) PRIMARY KEY,
				account_id varchar(36) NOT NULL,
				amount bigint NOT NULL,
				created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
	
				FOREIGN KEY (account_id) REFERENCES accounts(id)
			);
    `

	// At this point in the v2, this view query should be applied to the same db as the previous query.
	// For my purposes in this version thats fine, we can keep this since the volume will not be big and i have more things to do :)
	// If you want to test this I would suggest a few things
	// 1 - Setup one read replica [assuming you test it with the volume expected for this version 100/req/sec healthy paying shit machines]
	// 2 - Instead of view move it to a materialized view, that will provide you better read performance in exchange of some delay for updates and some work around updates
	// !In real world payments/banks application you for sure will have a way robust system for summary users/companies transactions so this is only for ref and play around
	viewQuery := `
		CREATE OR REPLACE VIEW account_payments_quarterly_summary AS
		WITH base AS (
			SELECT *
			FROM payments
			WHERE created_at >= CURRENT_DATE - INTERVAL '4 months'
		),
			 quarterly AS (
				 SELECT
					 base.from_account_id as account_id,
					 COUNT(*) as total_transactions,
					 ROUND(MAX(amount))::INTEGER AS max_ever_transaction_value,
					 ROUND(AVG(amount))::INTEGER AS avg_transaction_value
				 FROM base
				 GROUP BY account_id
			 ),
			 daily AS (
				 SELECT
					 base.from_account_id as account_id,
					 DATE(created_at) as day,
					 COUNT(*) as total_transactions_daily
				 FROM base
				 GROUP BY account_id, DATE(created_at)
			 ),
			 weekly AS (
				 SELECT
					 base.from_account_id as account_id,
					 DATE_TRUNC('week', created_at) as week_start,
					 COUNT(*) as total_transactions_weekly
				 FROM base
				 GROUP BY account_id, DATE_TRUNC('week', created_at)
			 )
		SELECT DISTINCT
			q.account_id,
			q.total_transactions,
			q.max_ever_transaction_value,
			q.avg_transaction_value,
			ROUND(AVG(d.total_transactions_daily))::INTEGER as avg_daily_transactions,
			ROUND(AVG(w.total_transactions_weekly))::INTEGER as avg_weekly_transactions
		FROM quarterly q
				 LEFT JOIN daily d ON q.account_id = d.account_id
				 LEFT JOIN weekly w ON q.account_id = w.account_id
		GROUP BY
			q.account_id,
			q.total_transactions,
			q.max_ever_transaction_value,
			q.avg_transaction_value
		ORDER BY q.account_id;
	`
	ctx := context.TODO()
	if _, err := db.Exec(ctx, query); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(ctx, viewQuery); err != nil {
		log.Fatal(err)
	}
}
