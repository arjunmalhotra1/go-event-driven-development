package http

import (
	"fmt"
	"net/http"
	"tickets/entities"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ticketsStatusRequest struct {
	Tickets []ticketStatusRequest `json:"tickets"`
}

type ticketStatusRequest struct {
	TicketID      string         `json:"ticket_id"`
	Status        string         `json:"status"`
	Price         entities.Money `json:"price"`
	CustomerEmail string         `json:"customer_email"`
	BookingID     string         `json:"booking_id"`
}

type postShowRequest struct {
	DeadNationId    string `json:"dead_nation_id"`
	NumberOfTickets int64  `json:"number_of_tickets"`
	StartTime       string `json:"start_time"`
	Title           string `json:"title"`
	Venue           string `json:"venue"`
}

func (h Handler) PostTicketsStatus(c echo.Context) error {
	var request ticketsStatusRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	for _, ticket := range request.Tickets {
		if ticket.Status == "confirmed" {
			var header entities.EventHeader
			idemKey := c.Request().Header.Get("Idempotency-Key")
			if idemKey != "" {
				header = entities.NewEventHeaderWithIdempotencyKey(idemKey + ticket.TicketID)
			} else {
				header = entities.NewEventHeader()
			}

			event := entities.TicketBookingConfirmed{
				Header:        header,
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			}

			err = h.eventBus.Publish(c.Request().Context(), event)
			if err != nil {
				return fmt.Errorf("failed to publish TicketBookingConfirmed event: %w", err)
			}
		} else if ticket.Status == "canceled" {
			var header entities.EventHeader
			idemKey := c.Request().Header.Get("Idempotency-Key")
			if idemKey != "" {
				header = entities.NewEventHeaderWithIdempotencyKey(idemKey + ticket.TicketID)
			} else {
				header = entities.NewEventHeader()
			}

			event := entities.TicketBookingCanceled{
				Header:        header,
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			}

			err = h.eventBus.Publish(c.Request().Context(), event)
			if err != nil {
				return fmt.Errorf("failed to publish TicketBookingCanceled event: %w", err)

			}
		} else {
			return fmt.Errorf("unknown ticket status: %s", ticket.Status)
		}
	}

	return c.NoContent(http.StatusOK)
}

func (h Handler) GetAllTickets(ctx echo.Context) error {
	res, err := h.ticketRepository.GetAll(ctx.Request().Context())
	if err != nil {
		return fmt.Errorf("failed to publish get all tickets event %w", err)
	}

	return ctx.JSON(http.StatusOK, res)
}

func (h Handler) PostShow(ctx echo.Context) error {

	var request postShowRequest
	err := ctx.Bind(&request)
	if err != nil {
		return err
	}

	showID := uuid.NewString()
	show := entities.Show{
		ShowId:          showID,
		DeadNationId:    request.DeadNationId,
		NumberOfTickets: request.NumberOfTickets,
		StartTime:       request.StartTime,
		Title:           request.Title,
		Venue:           request.Venue,
	}

	h.showRepository.Add(ctx.Request().Context(), show)

	response := struct {
		ShowID string `json:"show_id"`
	}{
		ShowID: showID,
	}
	return ctx.JSON(http.StatusCreated, response)
}
