package db

import (
	"context"
	"fmt"
	"tickets/entities"

	"github.com/jmoiron/sqlx"
)

type TicketsRepository struct {
	db *sqlx.DB
}

var insertQuery string = `Insert into tickets (ticket_id, price_amount, price_currency, customer_email)
VALUES (:ticket_id, :price.amount, :price.currency, :customer_email)`

var deleteQuery string = `Delete from tickets where ticket_id = ":ticket_id"`

func NewTicketsRepository(db *sqlx.DB) TicketsRepository {
	if db == nil {
		panic("Database client is nil")
	}

	return TicketsRepository{db: db}
}

func (t TicketsRepository) Add(ctx context.Context, ticket entities.Ticket) error {
	_, err := t.db.NamedExec(insertQuery, ticket)
	if err != nil {
		return fmt.Errorf("could not save the ticket: %w", err)
	}

	return nil
}

func (t TicketsRepository) Delete(ctx context.Context, ticketID string) error {
	_, err := t.db.ExecContext(
		ctx,
		`DELETE FROM tickets WHERE ticket_id = $1`,
		ticketID,
	)

	if err != nil {
		return fmt.Errorf("could not remove ticket: %w", err)
	}

	return nil
}
