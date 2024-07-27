package db

import (
	"context"
	"errors"
	"fmt"
	"tickets/entities"

	"github.com/jmoiron/sqlx"
)

type BookingRepository struct {
	db *sqlx.DB
}

var insertBookingQuery string = `
Insert into bookings (booking_id, show_id, number_of_tickets, customer_email) VALUES (:booking_id, :show_id, :number_of_tickets, :customer_email) ON CONFLICT DO NOTHING`

func NewBookingRepository(db *sqlx.DB) BookingRepository {
	return BookingRepository{
		db: db,
	}
}

func (br BookingRepository) Add(ctx context.Context, booking entities.Booking) error {
	tx, err := br.db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin the transaction: %w", err)
	}

	defer func() {
		if err != nil {
			rollBackErr := tx.Rollback()
			err = errors.Join(err, rollBackErr)
			return
		}
		err = tx.Commit()
	}()

	_, err = br.db.NamedExecContext(ctx, insertBookingQuery, booking)
	if err != nil {
		return fmt.Errorf("could not save the booking: %w", err)
	}

	return nil
}
