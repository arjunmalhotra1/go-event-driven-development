package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

var createTicketsTable string = `CREATE TABLE IF NOT EXISTS tickets (
	ticket_id UUID PRIMARY KEY,
	price_amount NUMERIC(10,2) NOT NULL,
	price_currency CHAR(3) NOT NULL,
	customer_email VARCHAR(255) NOT NULL
	);`

func InitializeDatabaseSchema(db *sqlx.DB) error {
	res, err := db.Exec(createTicketsTable)

	if err != nil {
		return fmt.Errorf("could not initialize database schema: %w", err)
	}

	fmt.Println(res)

	return nil
}
