package http

import (
	"context"
	"tickets/entities"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type Handler struct {
	eventBus              *cqrs.EventBus
	spreadsheetsAPIClient SpreadsheetsAPI
	ticketRepository      TicketRepository
}

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, spreadsheetName string, row []string) error
}

// NewHttpRouter in the service.go file creates an http handler.
// http handler calls the e.GET("/tickets", handler.GetAllTickets)
// handler.GetAllTickets calls the h.ticketRepository.GetAll(ctx.Request().Context())
// which is why the field "ticketRepository" in the http handler must have GetAll(ctx.Request().Context())
// Hence the "TicketRepository" is an interface in the
// type TicketRepository interface {
// GetAll(ctx context.Context) ([]entities.Ticket, error)

type TicketRepository interface {
	GetAll(ctx context.Context) ([]entities.Ticket, error)
}
