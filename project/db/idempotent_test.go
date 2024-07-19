package db_test

import (
	"context"
	"os"
	"sync"
	"testing"
	ticketsDb "tickets/db"
	"tickets/entities"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

var getDbOnce sync.Once
var db *sqlx.DB

func getDb() *sqlx.DB {
	getDbOnce.Do(func() {
		var err error
		db, err = sqlx.Open("postgres", os.Getenv("POSTGRES_URL"))
		if err != nil {
			panic(err)
		}
	})

	return db
}

func TestTicketRepository_Add_Idempotency(t *testing.T) {
	ctx := context.Background()
	db := getDb()
	err := ticketsDb.InitializeDatabaseSchema(db)
	require.NoError(t, err)

	repo := ticketsDb.NewTicketsRepository(db)

	ticketToAdd := entities.Ticket{
		TicketID: uuid.NewString(),
		Price: entities.Money{
			Amount:   "1234",
			Currency: "dol",
		},
		CustomerEmail: "abcd@xyz.com",
	}

	for i := 0; i < 2; i++ {
		err = repo.Add(ctx, ticketToAdd)
		require.NoError(t, err)

		// probably it would be good to have a method to get ticket by ID
		tickets, err := repo.GetAll(ctx)
		require.NoError(t, err)
		foundTickets := lo.Filter(tickets, func(ticket entities.Ticket, _ int) bool {
			return ticket.TicketID == ticketToAdd.TicketID
		})

		// add should be idempotent, so the method should always return 1
		require.Len(t, foundTickets, 1)
	}

}
