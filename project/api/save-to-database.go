package api

import (
	"context"
	"tickets/db"

	"github.com/jmoiron/sqlx"
)

type DatabaseClient struct {
	dbConn *sqlx.DB
}

var insertQuery string = "Insert into tickets (ticket.ID, ticket."

func NewDatabaseClient(dbConn *sqlx.DB) *DatabaseClient {
	if dbConn == nil {
		panic("Database client is nil")
	}

	return &DatabaseClient{dbConn: dbConn}
}

func (db DatabaseClient) SaveTicketToDatabase(ctx context.Context, ticket db.Ticket) error {
	_, err := db.dbConn.Exec(insertQuery)
	if err != nil {
		panic(err)
	}

	return nil
}
