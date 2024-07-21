package api

import (
	"context"
	"fmt"
	"net/http"
	"tickets/entities"

	"github.com/ThreeDotsLabs/go-event-driven/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/common/clients/receipts"
)

type ReceiptsClient struct {
	clients *clients.Clients
}

func NewReceiptsClient(clients *clients.Clients) *ReceiptsClient {
	return &ReceiptsClient{
		clients: clients,
	}
}

func (c ReceiptsClient) IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) (entities.IssueReceiptResponse, error) {

	body := receipts.CreateReceipt{
		IdempotencyKey: &request.IdempotencyKey,
		TicketId:       request.TicketID,
		Price: receipts.Money{
			MoneyAmount:   request.Price.Amount,
			MoneyCurrency: request.Price.Currency,
		},
	}

	receiptsResp, err := c.clients.Receipts.PutReceiptsWithResponse(ctx, body)
	if err != nil {
		return entities.IssueReceiptResponse{}, err
	}

	switch receiptsResp.StatusCode() {
	case http.StatusOK:
		// receipt already exists
		return entities.IssueReceiptResponse{
			ReceiptNumber: receiptsResp.JSON200.Number,
			IssuedAt:      receiptsResp.JSON200.IssuedAt,
		}, nil
	case http.StatusCreated:
		// receipt was created
		return entities.IssueReceiptResponse{
			ReceiptNumber: receiptsResp.JSON201.Number,
			IssuedAt:      receiptsResp.JSON201.IssuedAt,
		}, nil
	default:
		return entities.IssueReceiptResponse{}, fmt.Errorf("unexpected status code: %v", receiptsResp.StatusCode())
	}
}
