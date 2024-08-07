package event

import (
	"context"
	"fmt"
	"tickets/entities"
	"time"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
)

func (h Handler) PrintTicket(ctx context.Context, event *entities.TicketBookingConfirmed) error {
	log.FromContext(ctx).Info("Printing ticket")

	ticketHTML := `
		<html>
			<head>
				<title>Ticket</title>
			</head>
			<body>
				<h1>Ticket ` + event.TicketID + `</h1>
				<p>Price: ` + event.Price.Amount + ` ` + event.Price.Currency + `</p>	
			</body>
		</html>
`

	ticketFile := event.TicketID + "-ticket.html"

	err := h.filesAPI.UploadTicket(ctx, ticketFile, ticketHTML)
	if err != nil {
		return fmt.Errorf("failed to upload the ticket file %w", err)
	}

	ticketPrintedEvent := entities.TicketPrinted{
		Header: entities.EventHeader{
			ID:          event.TicketID,
			PublishedAt: time.Now(),
		},
		TicketID: event.TicketID,
		FileName: ticketFile,
	}

	err = h.eventBus.Publish(ctx, ticketPrintedEvent)
	if err != nil {
		return fmt.Errorf("failed to publish TicketBookingConfirmed event: %w", err)
	}

	return nil
}
