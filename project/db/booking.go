package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"tickets/entities"
	"tickets/message"
	"tickets/message/event"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"

	watermillSQL "github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/components/forwarder"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

var outboxTopic = "events_to_forward"

type BookingRepository struct {
	db *sqlx.DB
}

var insertBookingQuery string = `
Insert into bookings (booking_id, show_id, number_of_tickets, customer_email) VALUES (:booking_id, :show_id, :number_of_tickets, :customer_email) ON CONFLICT DO NOTHING`

func NewBookingRepository(db *sqlx.DB) BookingRepository {
	return BookingRepository{
		db: db,
	}
}

func (br BookingRepository) Add(ctx context.Context, booking entities.Booking) error {
	tx, err := br.db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin the transaction: %w", err)
	}

	logger := log.NewWatermill(log.FromContext(context.Background()))

	publisher, err := watermillSQL.NewPublisher(tx,
		watermillSQL.PublisherConfig{
			SchemaAdapter: watermillSQL.DefaultPostgreSQLSchema{},
		},
		logger,
	)
	if err != nil {
		return fmt.Errorf("error creating new publisher with watermillSQL")
	}

	forwarderPublsiher := forwarder.NewPublisher(publisher, forwarder.PublisherConfig{
		ForwarderTopic: outboxTopic,
	})

	event := event.ForwarderEvent{
		BookingId:       booking.BookingID,
		NumberOfTickets: booking.NumberOfTickets,
		CustomerEmail:   booking.CustomerEmail,
		ShowId:          booking.ShowID,
	}
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error marshalling the msg")
	}

	msg := watermillMessage.NewMessage(watermill.NewUUID(), eventBytes)

	err = forwarderPublsiher.Publish("BookingMade", msg)
	if err != nil {
		return fmt.Errorf("error publishing msg to the BookingMade topic")
	}

	redisClient := message.NewRedisClient(os.Getenv("REDIS_ADDR"))
	defer redisClient.Close()
	err = RunForwarder(br.db, redisClient, outboxTopic, logger)

	defer func() {
		if err != nil {
			rollBackErr := tx.Rollback()
			err = errors.Join(err, rollBackErr)
			return
		}
		err = tx.Commit()
	}()

	_, err = br.db.NamedExecContext(ctx, insertBookingQuery, booking)
	if err != nil {
		return fmt.Errorf("could not save the booking: %w", err)
	}

	return nil
}

func RunForwarder(
	db *sqlx.DB,
	rdb *redis.Client,
	outboxTopic string,
	logger watermill.LoggerAdapter,
) error {

	sqlSubscriberConfig := watermillSQL.SubscriberConfig{
		SchemaAdapter:  watermillSQL.DefaultPostgreSQLSchema{},
		OffsetsAdapter: watermillSQL.DefaultPostgreSQLOffsetsAdapter{},
	}
	sqlSubscriber, err := watermillSQL.NewSubscriber(db, sqlSubscriberConfig, logger)

	if err != nil {
		return fmt.Errorf("error creating new subscriber from sql %w", err)
	}

	sqlSubscriber.SubscribeInitialize(outboxTopic)

	redisPubConfig := redisstream.PublisherConfig{Client: rdb}
	redisPub, err := redisstream.NewPublisher(redisPubConfig, logger)

	if err != nil {
		return fmt.Errorf("error creating the redis publisher")
	}

	fwdr, err := forwarder.NewForwarder(sqlSubscriber, redisPub, logger, forwarder.Config{
		ForwarderTopic: outboxTopic,
	})

	if err != nil {
		return fmt.Errorf("error creating the new forwarder %w", err)
	}
	go func() {
		err = fwdr.Run(context.Background())
		if err != nil {
			panic(err)
		}
	}()

	<-fwdr.Running()

	return nil
}
