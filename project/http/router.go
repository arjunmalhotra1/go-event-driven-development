package http

import (
	"net/http"

	libHttp "github.com/ThreeDotsLabs/go-event-driven/common/http"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/labstack/echo/v4"
)

func NewHttpRouter(eventBus *cqrs.EventBus, spreadsheetsAPIClient SpreadsheetsAPI, ticketRepository TicketRepository, showRepository ShowRepository) *echo.Echo {
	e := libHttp.NewEcho()

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	handler := Handler{
		eventBus:              eventBus,
		spreadsheetsAPIClient: spreadsheetsAPIClient,
		ticketRepository:      ticketRepository,
		showRepository:        showRepository,
	}

	e.POST("/tickets-status", handler.PostTicketsStatus)

	e.POST("/shows", handler.PostShow)

	e.GET("/tickets", handler.GetAllTickets)
	return e
}
