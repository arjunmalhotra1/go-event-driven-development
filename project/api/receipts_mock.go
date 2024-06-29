package api

import (
	"context"
	"sync"
	"tickets/entities"
	"time"
)

type ReceiptsServiceMock struct {
	IssuedReceipts []entities.IssueReceiptRequest
	Mock           sync.Mutex
}

func (r *ReceiptsServiceMock) IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) (entities.IssueReceiptResponse, error) {
	var response entities.IssueReceiptResponse
	r.Mock.Lock()
	defer r.Mock.Unlock()
	r.IssuedReceipts = append(r.IssuedReceipts, request)

	response.ReceiptNumber = "1234"
	response.IssuedAt = time.Now()
	return response, nil
}
