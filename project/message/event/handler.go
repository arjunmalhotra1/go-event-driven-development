package event

import (
	"context"
	"tickets/db"
	"tickets/entities"
)

type Handler struct {
	spreadsheetsService         SpreadsheetsAPI
	receiptsService             ReceiptsService
	saveTicketToDatabaseService SaveTicketToDatabaseService
}

func NewHandler(
	spreadsheetsService SpreadsheetsAPI,
	receiptsService ReceiptsService,
	saveTicketToDatabaseService SaveTicketToDatabaseService,
) Handler {
	if spreadsheetsService == nil {
		panic("missing spreadsheetsService")
	}
	if receiptsService == nil {
		panic("missing receiptsService")
	}

	return Handler{
		spreadsheetsService:         spreadsheetsService,
		receiptsService:             receiptsService,
		saveTicketToDatabaseService: saveTicketToDatabaseService,
	}
}

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) (entities.IssueReceiptResponse, error)
}

type SaveTicketToDatabaseService interface {
	SaveTicketToDatabase(ctx context.Context, ticket db.Ticket) error
}
