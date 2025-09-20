package repo

import (
	"context"

	"github.com/cprakhar/uber-clone/services/payment-service/types"
)

type Service interface {
	CreatePaymentSession(ctx context.Context, tripID, riderID, driverID string, amount int64, currency string) (*types.PaymentIntent, error)
}

type PaymentProcessor interface {
	CreatePaymentSession(ctx context.Context, amount int64, currency string, metadata map[string]string) (string, error)
}
