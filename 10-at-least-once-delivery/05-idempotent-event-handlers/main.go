package main

import (
	"context"
	"fmt"
)

type PaymentTaken struct {
	PaymentID string
	Amount    int
}

var paymentMap map[string]struct{}

type PaymentsHandler struct {
	repo *PaymentsRepository
}

func NewPaymentsHandler(repo *PaymentsRepository) *PaymentsHandler {
	return &PaymentsHandler{repo: repo}
}

func (p *PaymentsHandler) HandlePaymentTaken(ctx context.Context, event *PaymentTaken) error {
	return p.repo.SavePaymentTaken(ctx, event)
}

type PaymentsRepository struct {
	payments []PaymentTaken
}

func (p *PaymentsRepository) Payments() []PaymentTaken {
	return p.payments
}

func NewPaymentsRepository() *PaymentsRepository {
	paymentMap = make(map[string]struct{})
	return &PaymentsRepository{}
}

func (p *PaymentsRepository) SavePaymentTaken(ctx context.Context, event *PaymentTaken) error {
	fmt.Println(event.PaymentID)
	if _, ok := paymentMap[event.PaymentID]; !ok {
		p.payments = append(p.payments, *event)
		paymentMap[event.PaymentID] = struct{}{}
	}

	return nil
}
