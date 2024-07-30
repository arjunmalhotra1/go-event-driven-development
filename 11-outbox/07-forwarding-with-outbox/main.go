package main

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/components/forwarder"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	_ "github.com/lib/pq"
)

func RunForwarder(
	db *sqlx.DB,
	rdb *redis.Client,
	outboxTopic string,
	logger watermill.LoggerAdapter,
) error {

	sqlSubscriberConfig := sql.SubscriberConfig{
		SchemaAdapter:  sql.DefaultPostgreSQLSchema{},
		OffsetsAdapter: sql.DefaultPostgreSQLOffsetsAdapter{},
	}
	sqlSubscriber, err := sql.NewSubscriber(db, sqlSubscriberConfig, logger)

	if err != nil {
		return fmt.Errorf("error creating new subscriber from sql %w", err)
	}

	sqlSubscriber.SubscribeInitialize(outboxTopic)

	// o, err := sqlSubscriber.Subscribe(context.Background(), outboxTopic)
	// if err != nil {
	// 	return fmt.Errorf("error subscribing from the sqlSubscriber %w", err)
	// }

	// tx, err := db.BeginTx(context.Background(), nil)
	// if err != nil {
	// 	return fmt.Errorf("error starting the transaction %w", err)
	// }

	// publisher, err := sql.NewPublisher(
	// 	tx,
	// 	sql.PublisherConfig{
	// 		SchemaAdapter: sql.DefaultPostgreSQLSchema{},
	// 	},
	// 	logger,
	// )

	// if err != nil {
	// 	return fmt.Errorf("error creating a publisher %w", err)
	// }

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
