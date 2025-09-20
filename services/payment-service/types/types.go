package types

import "time"

// PaymentStatus represents the current status of a payment
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSuccess   PaymentStatus = "success"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

// Payment represents a payment transaction
type Payment struct {
	ID              string        `json:"id"`
	TripID          string        `json:"tripID"`
	RiderID         string        `json:"riderID"`
	Amount          int64         `json:"amount"`   // Amount in cents
	Currency        string        `json:"currency"` // e.g., "usd"
	Status          PaymentStatus `json:"status"`
	StripeSessionID string        `json:"stripeSessionID"`
	CreatedAt       time.Time     `json:"createdAt"`
	UpdatedAt       time.Time     `json:"updatedAt"`
}

// PaymentIntent represents the intent to collect a payment
type PaymentIntent struct {
	ID              string    `json:"id"`
	TripID          string    `json:"tripID"`
	RiderID         string    `json:"riderID"`
	DriverID        string    `json:"driverID"`
	Amount          int64     `json:"amount"`
	Currency        string    `json:"currency"`
	StripeSessionID string    `json:"stripeSessionID"`
	CreatedAt       time.Time `json:"createdAt"`
}

// PaymentConfig holds the configuration for the payment service
type PaymentConfig struct {
	StripeSecretKey     string `json:"stripeSecretKey"`
	StripeWebhookSecret string `json:"stripeWebhookSecret"`
	Currency            string `json:"currency"`
	SuccessURL          string `json:"successURL"`
	CancelURL           string `json:"cancelURL"`
}
