package events

import (
	"context"
	"encoding/json"
	"log"
	"time"

	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/cprakhar/uber-clone/services/payment-service/repo"
	"github.com/cprakhar/uber-clone/shared/contracts"
	"github.com/cprakhar/uber-clone/shared/messaging"
	"github.com/cprakhar/uber-clone/shared/messaging/kafka"
)

type TripConsumer struct {
	kfClient *kafka.KafkaClient
	svc      repo.Service
}

func NewTripConsumer(kfClient *kafka.KafkaClient, svc repo.Service) *TripConsumer {
	return &TripConsumer{kfClient: kfClient, svc: svc}
}

func (tc *TripConsumer) Consume(ctx context.Context, topics []string) error {
	return tc.kfClient.Consumer.SubscribeAndConsume(ctx, topics,
		func(ctx context.Context, m *ckafka.Message) error {
			var kafkaMsg contracts.KafkaMessage
			if err := json.Unmarshal(m.Value, &kafkaMsg); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
			}

			var payload messaging.PaymentTripResponseData
			if kafkaMsg.Data != nil {
				if err := json.Unmarshal(kafkaMsg.Data, &payload); err != nil {
					log.Printf("Failed to unmarshal payload: %v", err)
				}
			}

			switch *m.TopicPartition.Topic {
			case contracts.PaymentCmdCreateSession:
				if err := tc.handleTripAccepted(ctx, payload); err != nil {
					log.Printf("Failed to handle trip accepted: %v", err)
					return err
				}
			}
			return nil
		},
	)
}

func (tc *TripConsumer) handleTripAccepted(ctx context.Context, payload messaging.PaymentTripResponseData) error {
	log.Printf("Processing payment for trip %s with amount %.2f", payload.TripID, payload.Amount)

	paymentSession, err := tc.svc.CreatePaymentSession(ctx,
		payload.TripID,
		payload.RiderID,
		payload.DriverID,
		int64(payload.Amount),
		payload.Currency,
	)

	if err != nil {
		log.Printf("Failed to create payment session: %v", err)
		return err
	}

	log.Printf("Payment session created: %s", paymentSession.StripeSessionID)

	paymentPayload := messaging.PaymentEventSessionCreatedData{
		TripID:    payload.TripID,
		SessionID: paymentSession.StripeSessionID,
		Amount:    float64(payload.Amount) / 100, // converting paise to rupees
		Currency:  payload.Currency,
	}

	data, err := json.Marshal(paymentPayload)
	if err != nil {
		log.Printf("Failed to marshal payload: %v", err)
		return err
	}

	if err := tc.kfClient.Producer.SendMessageAndWait(ctx, contracts.PaymentEventSessionCreated,
		&contracts.KafkaMessage{
			EntityID: payload.RiderID,
			Data:     data,
		},
		30*time.Second,
	); err != nil {
		log.Printf("Failed to send payment session created message: %v", err)
		return err
	}

	log.Printf("Payment session created message sent for trip %s", payload.TripID)
	return nil
}
