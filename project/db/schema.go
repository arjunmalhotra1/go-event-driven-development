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

var createShowsTable string = `CREATE TABLE IF NOT EXISTS shows (
	show_id UUID PRIMARY KEY,
	dead_nation_id UUID NOT NULL,
	number_of_tickets INT NOT NULL,
	start_time TIMESTAMP NOT NULL,
	title VARCHAR(255) NOT NULL,
	venue VARCHAR(255) NOT NULL,

	UNIQUE (dead_nation_id)
	);`

var createBookingTable string = `CREATE TABLE IF NOT EXISTS bookings (
	booking_id UUID PRIMARY KEY,
	show_id UUID NOT NULL,
	number_of_tickets INT NOT NULL,
	customer_email VARCHAR(255) NOT NULL,
	FOREIGN KEY (show_id) REFERENCES shows(show_id)
	);`

func InitializeDatabaseSchema(db *sqlx.DB) error {
	_, err := db.Exec(createTicketsTable)
	if err != nil {
		return fmt.Errorf("could not initialize createTicketsTable database schema: %w", err)
	}

	_, err = db.Exec(createShowsTable)
	if err != nil {
		return fmt.Errorf("could not initialize createShowsTable database schema : %w", err)
	}

	_, err = db.Exec(createBookingTable)
	if err != nil {
		return fmt.Errorf("could not initialize createBookingsTable database schema : %w", err)
	}

	return nil
}
