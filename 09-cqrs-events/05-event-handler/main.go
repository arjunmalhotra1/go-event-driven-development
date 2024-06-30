package main

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type FollowRequestSent struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type EventsCounter interface {
	CountEvent() error
}

type EventsCounterHandler struct {
	eventCounter EventsCounter
}

func (ech EventsCounterHandler) NotifcationEvent(ctx context.Context, event *FollowRequestSent) error {
	return ech.eventCounter.CountEvent()
}

func NewFollowRequestSentHandler(counter EventsCounter) cqrs.EventHandler {
	h := EventsCounterHandler{
		eventCounter: counter,
	}
	return cqrs.NewEventHandler("", h.NotifcationEvent)
}
