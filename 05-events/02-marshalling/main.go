package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

type PaymentCompleted struct {
	PaymentID   string `json:"payment_id"`
	OrderID     string `json:"order_id"`
	CompletedAt string `json:"completed_at"`
}

type OrderCompleted struct {
	OrderID     string `json:"order_id"`
	ConfirmedAt string `json:"confirmed_at"`
}

func main() {
	logger := watermill.NewStdLogger(false, false)

	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	sub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		panic(err)
	}

	pub, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		panic(err)
	}

	router.AddHandler("test-marshall", "payment-completed", sub, "order-confirmed", pub,
		func(msg *message.Message) ([]*message.Message, error) {
			var pCompletedEvent PaymentCompleted
			var oCompletedEvent OrderCompleted

			err := json.Unmarshal(msg.Payload, &pCompletedEvent)
			if err != nil {
				panic("failed to unmarshall")
			}

			oCompletedEvent.ConfirmedAt = pCompletedEvent.CompletedAt
			oCompletedEvent.OrderID = pCompletedEvent.OrderID

			fmt.Println("o:", oCompletedEvent.ConfirmedAt)

			oCompletedBytes, err := json.Marshal(oCompletedEvent)
			if err != nil {
				panic("failed to marshal")
			}

			newMsg := message.NewMessage(watermill.NewUUID(), oCompletedBytes)

			return []*message.Message{newMsg}, nil
		})

	err = router.Run(context.Background())
	if err != nil {
		panic(err)
	}
}
