package event

import (
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

func NewEventBus(pub message.Publisher) *cqrs.EventBus {
	config := cqrs.EventBusConfig{
		GeneratePublishTopic: func(params cqrs.GenerateEventPublishTopicParams) (string, error) {
			return params.EventName, nil
		},
		Marshaler: cqrs.JSONMarshaler{GenerateName: cqrs.StructName},
	}

	eventBus, err := cqrs.NewEventBusWithConfig(pub, config)
	if err != nil {
		panic(err)
	}

	return eventBus
}
