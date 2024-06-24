package main

import (
	"context"
	"sync"
	"time"
)

type IssueReceiptRequest struct {
	TicketID string `json:"ticket_id"`
	Price    Money  `json:"price"`
}

type Money struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type IssueReceiptResponse struct {
	ReceiptNumber string    `json:"number"`
	IssuedAt      time.Time `json:"issued_at"`
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request IssueReceiptRequest) (IssueReceiptResponse, error)
}

type ReceiptsServiceMock struct {
	IssuedReceipts []IssueReceiptRequest
	mock           sync.Mutex
}

func (r *ReceiptsServiceMock) IssueReceipt(ctx context.Context, request IssueReceiptRequest) (IssueReceiptResponse, error) {
	var response IssueReceiptResponse
	r.mock.Lock()
	defer r.mock.Unlock()
	r.IssuedReceipts = append(r.IssuedReceipts, request)

	response.ReceiptNumber = "1234"
	response.IssuedAt = time.Now()
	return response, nil
}
