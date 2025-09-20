package service

import (
	"context"
	"fmt"
	"time"

	"github.com/cprakhar/uber-clone/services/payment-service/repo"
	"github.com/cprakhar/uber-clone/services/payment-service/types"
	"github.com/google/uuid"
)

type paymentService struct {
	paymentProcessor repo.PaymentProcessor
}

func NewPaymentService(r repo.PaymentProcessor) repo.Service {
	return &paymentService{paymentProcessor: r}
}

func (s *paymentService) CreatePaymentSession(
	ctx context.Context,
	tripID, riderID, driverID string,
	amount int64,
	currency string) (*types.PaymentIntent, error) {

	metadata := map[string]string{
		"tripID":   tripID,
		"riderID":  riderID,
		"driverID": driverID,
	}

	sessionID, err := s.paymentProcessor.CreatePaymentSession(ctx, amount, currency, metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment session: %w", err)
	}

	return &types.PaymentIntent{
		ID:              uuid.New().String(),
		TripID:          tripID,
		RiderID:         riderID,
		DriverID:        driverID,
		Amount:          amount,
		Currency:        currency,
		StripeSessionID: sessionID,
		CreatedAt:       time.Now(),
	}, nil

}
