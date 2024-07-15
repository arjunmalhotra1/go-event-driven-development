package event

import (
	"context"
	"fmt"
	"tickets/db"
	"tickets/entities"
)

func (h Handler) SaveTicketToDatabase(ctx context.Context, ticket *entities.Ticket) error {

	t := db.Ticket{
		TicketID:      ticket.TicketID,
		PriceAmount:   ticket.Price.Amount,
		PriceCurrency: ticket.Price.Currency,
		CustomerEmail: ticket.CustomerEmail,
	}

	err := h.saveTicketToDatabaseService.SaveTicketToDatabase(ctx, t)

	if err != nil {
		return fmt.Errorf("failed to issue receipt: %w", err)
	}

	return nil
}
