package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ThreeDotsLabs/go-event-driven/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/common/clients/spreadsheets"
)

type SpreadsheetsClient struct {
	clients *clients.Clients
}

func NewSpreadsheetsAPIClient(clients *clients.Clients) *SpreadsheetsClient {
	if clients == nil {
		panic("NewSpreadsheetsAPIClient: clients is nil")
	}

	return &SpreadsheetsClient{
		clients: clients,
	}
}

func (c SpreadsheetsClient) AppendRow(ctx context.Context, spreadsheetName string, row []string) error {
	request := spreadsheets.PostSheetsSheetRowsJSONRequestBody{
		Columns: row,
	}

	sheetsResp, err := c.clients.Spreadsheets.PostSheetsSheetRowsWithResponse(ctx, spreadsheetName, request)
	if err != nil {
		return fmt.Errorf("failed to post row: %w", err)
	}

	if sheetsResp.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed to post row: unexpected status code %d", sheetsResp.StatusCode())
	}

	return nil
}
