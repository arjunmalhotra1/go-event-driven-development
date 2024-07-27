package http

import (
	"fmt"
	"net/http"
	"tickets/entities"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h Handler) PostBooking(ctx echo.Context) error {
	var booking entities.Booking
	booking.BookingID = uuid.New()

	if err := ctx.Bind(&booking); err != nil {
		return fmt.Errorf("error binding the request")
	}

	if booking.NumberOfTickets < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, "number of tickets must be greater than 0")
	}

	h.bookingRepository.Add(ctx.Request().Context(), booking)

	return ctx.JSON(http.StatusCreated, booking)
}
