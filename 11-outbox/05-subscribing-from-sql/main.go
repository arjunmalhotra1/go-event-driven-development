package main

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func SubscribeForMessages(db *sqlx.DB, topic string, logger watermill.LoggerAdapter) (<-chan *message.Message, error) {

	sqlSConfig := sql.SubscriberConfig{
		SchemaAdapter:  sql.DefaultPostgreSQLSchema{},
		OffsetsAdapter: sql.DefaultPostgreSQLOffsetsAdapter{},
		//InitializeSchema: true,
	}

	sqlSubscriber, err := sql.NewSubscriber(db, sqlSConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("error creating new subscriber from sql %w", err)
	}

	sqlSubscriber.SubscribeInitialize(topic)

	o, err := sqlSubscriber.Subscribe(context.Background(), topic)
	if err != nil {
		return nil, fmt.Errorf("error subscribing from the sqlSubscriber %w", err)
	}
	return o, nil
}
