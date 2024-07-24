package db

import (
	"context"
	"fmt"

	"tickets/entities"

	"github.com/jmoiron/sqlx"
)

type ShowRepository struct {
	db *sqlx.DB
}

var insertShowQuery string = `Insert into shows (show_id, dead_nation_id, number_of_tickets, start_time, title, venue)
VALUES (:show_id, :dead_nation_id, :number_of_tickets, :start_time, :title, :venue) ON CONFLICT DO NOTHING`

func NewShowRepository(db *sqlx.DB) ShowRepository {
	if db == nil {
		panic("Database client is nil")
	}

	return ShowRepository{db: db}
}

func (sr ShowRepository) Add(ctx context.Context, show entities.Show) error {
	fmt.Println("Inside the show add function")
	_, err := sr.db.NamedExec(insertShowQuery, show)
	if err != nil {
		return fmt.Errorf("could not save the show: %w", err)
	}

	return nil
}
