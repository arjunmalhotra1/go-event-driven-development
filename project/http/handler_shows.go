package http

import (
	"fmt"
	"net/http"
	"tickets/entities"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h Handler) PostShow(ctx echo.Context) error {

	show := entities.Show{}
	if err := ctx.Bind(&show); err != nil {
		return err
	}

	show.ShowId = uuid.New()

	if err := h.showRepository.Add(ctx.Request().Context(), show); err != nil {
		return fmt.Errorf("failed to add show:%w", err)
	}

	return ctx.JSON(http.StatusCreated, show)
}
