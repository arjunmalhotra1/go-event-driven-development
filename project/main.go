package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ThreeDotsLabs/go-event-driven/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/common/clients/receipts"
	"github.com/ThreeDotsLabs/go-event-driven/common/clients/spreadsheets"
	commonHTTP "github.com/ThreeDotsLabs/go-event-driven/common/http"
	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type TicketsConfirmationRequest struct {
	Tickets []Ticket `json:"tickets"`
}

type Ticket struct {
	TicketID      string `json:"ticket_id"`
	Status        string `json:"status"`
	CustomerEmail string `json:"customer_email"`
	Price         Money  `json:"price"`
}

type EventHeader struct {
	ID          string `json:"id"`
	PublishedAt string `json:"published_at"`
}

type TicketBookingConfirmed struct {
	Header        EventHeader `json:"header"`
	TicketID      string      `json:"ticket_id"`
	CustomerEmail string      `json:"customer_email"`
	Price         Money       `json:"price"`
}

type TicketBookingCanceled struct {
	Header        EventHeader `json:"header"`
	TicketID      string      `json:"ticket_id"`
	CustomerEmail string      `json:"customer_email"`
	Price         Money       `json:"price"`
}

type Money struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type IssueReceiptRequest struct {
	TicketID string `json:"ticket_id"`
	Price    Money  `json:"price"`
}

type AppendToTrackerRequest struct {
	TicketID      string `json:"ticket_id"`
	CustomerEmail string `json:"customer_email"`
	Price         Money  `json:"price"`
}

func NewHeader() EventHeader {
	return EventHeader{
		ID:          uuid.NewString(),
		PublishedAt: time.Now().Format(time.RFC3339),
	}
}

func main() {
	log.Init(logrus.InfoLevel)

	clients, err := clients.NewClients(os.Getenv("GATEWAY_ADDR"), nil)
	if err != nil {
		panic(err)
	}

	spreadsheetsClient := NewSpreadsheetsClient(clients)
	receiptsClient := NewReceiptsClient(clients)

	watermillLogger := log.NewWatermill(logrus.NewEntry(logrus.StandardLogger()))

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	publisher, err := redisstream.NewPublisher(redisstream.PublisherConfig{Client: rdb},
		watermillLogger,
	)
	if err != nil {
		panic(err)
	}

	issueReceiptSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "issue-receipt",
	}, watermillLogger)
	if err != nil {
		panic(err)
	}

	appendToTrackerSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "append-to-tracker",
	}, watermillLogger)
	if err != nil {
		panic(err)
	}

	e := commonHTTP.NewEcho()
	e.POST("/tickets-status", func(c echo.Context) error {
		var request TicketsConfirmationRequest
		err := c.Bind(&request)
		if err != nil {
			return err
		}

		for _, ticket := range request.Tickets {

			if ticket.Status == "canceled" {
				event := TicketBookingCanceled{
					Header:        NewHeader(),
					TicketID:      ticket.TicketID,
					CustomerEmail: ticket.CustomerEmail,
					Price:         ticket.Price,
				}

				payload, err := json.Marshal(event)
				if err != nil {
					return err
				}

				msg := message.NewMessage(watermill.NewUUID(), payload)

				err = publisher.Publish("TicketBookingCanceled", msg)
				if err != nil {
					panic(err)
				}
			}

			if ticket.Status == "confirmed" {
				event := TicketBookingConfirmed{
					Header:        NewHeader(),
					TicketID:      ticket.TicketID,
					CustomerEmail: ticket.CustomerEmail,
					Price:         ticket.Price,
				}

				payload, err := json.Marshal(event)
				if err != nil {
					return err
				}

				msg := message.NewMessage(watermill.NewUUID(), payload)

				err = publisher.Publish("TicketBookingConfirmed", msg)
				if err != nil {
					panic(err)
				}
			}
		}

		return c.NoContent(http.StatusOK)
	})

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	router, err := message.NewRouter(message.RouterConfig{}, watermillLogger)
	if err != nil {
		panic(err)
	}

	router.AddMiddleware(SaveLogsMiddleware())

	router.AddNoPublisherHandler("print_ticket", "TicketBookingConfirmed", appendToTrackerSub, func(msg *message.Message) error {
		var event TicketBookingConfirmed
		err := json.Unmarshal(msg.Payload, &event)
		if err != nil {
			panic(err)
		}

		var payload AppendToTrackerRequest
		payload.TicketID = event.TicketID
		payload.CustomerEmail = event.CustomerEmail
		payload.Price = event.Price

		return spreadsheetsClient.AppendRow(context.Background(), "tickets-to-print", []string{payload.TicketID, payload.CustomerEmail, payload.Price.Amount, payload.Price.Currency})
	})

	router.AddNoPublisherHandler("issue_receipt", "TicketBookingConfirmed", issueReceiptSub, func(msg *message.Message) error {
		var event TicketBookingConfirmed
		err := json.Unmarshal(msg.Payload, &event)
		if err != nil {
			panic(err)
		}

		var issueReceiptReq IssueReceiptRequest
		issueReceiptReq.TicketID = event.TicketID
		issueReceiptReq.Price = event.Price

		return receiptsClient.IssueReceipt(context.Background(), issueReceiptReq)
	})

	router.AddNoPublisherHandler("cancel_ticket", "TicketBookingCanceled", appendToTrackerSub, func(msg *message.Message) error {
		var event TicketBookingConfirmed
		err := json.Unmarshal(msg.Payload, &event)
		if err != nil {
			panic(err)
		}

		var payload AppendToTrackerRequest
		payload.TicketID = event.TicketID
		payload.CustomerEmail = event.CustomerEmail
		payload.Price = event.Price

		return spreadsheetsClient.AppendRow(context.Background(), "tickets-to-refund", []string{payload.TicketID, payload.CustomerEmail, payload.Price.Amount, payload.Price.Currency})
	})

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return router.Run(context.Background())
	})

	logrus.Info("Server starting...")

	g.Go(func() error {

		// we don't want to start HTTP server before Watermill router (so service won't be healthy before it's ready)
		<-router.Running()

		err := e.Start(":8080")
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	g.Go(func() error {
		// Shutdown the HTTP server
		<-ctx.Done()
		return e.Shutdown(ctx)

	})

	// Will block until all goroutines finish
	err = g.Wait()
	if err != nil {
		panic(err)
	}

}

func SaveLogsMiddleware() func(message.HandlerFunc) message.HandlerFunc {
	return func(f message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			logrus.WithField("message_uuid", msg.UUID).Info("Handling a message")
			return f(msg)
		}
	}
}

type ReceiptsClient struct {
	clients *clients.Clients
}

func NewReceiptsClient(clients *clients.Clients) ReceiptsClient {
	return ReceiptsClient{
		clients: clients,
	}
}

func (c ReceiptsClient) IssueReceipt(ctx context.Context, request IssueReceiptRequest) error {
	body := receipts.PutReceiptsJSONRequestBody{
		TicketId: request.TicketID,
		Price: receipts.Money{
			MoneyAmount:   request.Price.Amount,
			MoneyCurrency: request.Price.Currency,
		},
	}

	receiptsResp, err := c.clients.Receipts.PutReceiptsWithResponse(ctx, body)
	if err != nil {
		return err
	}
	if receiptsResp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %v", receiptsResp.StatusCode())
	}

	return nil
}

type SpreadsheetsClient struct {
	clients *clients.Clients
}

func NewSpreadsheetsClient(clients *clients.Clients) SpreadsheetsClient {
	return SpreadsheetsClient{
		clients: clients,
	}
}

func (c SpreadsheetsClient) AppendRow(ctx context.Context, spreadsheetName string, row []string) error {
	request := spreadsheets.PostSheetsSheetRowsJSONRequestBody{
		Columns: row,
	}

	sheetsResp, err := c.clients.Spreadsheets.PostSheetsSheetRowsWithResponse(ctx, spreadsheetName, request)
	if err != nil {
		return err
	}
	if sheetsResp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %v", sheetsResp.StatusCode())
	}

	return nil
}
