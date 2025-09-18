package events

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/cprakhar/uber-clone/shared/messaging"
)

type TripEventConsumer struct {
	Kafka *messaging.Consumer
}

// NewTripEventConsumer creates a new TripEventConsumer with the given Kafka consumer.
func NewTripEventConsumer(kafka *messaging.Consumer) *TripEventConsumer {
	return &TripEventConsumer{Kafka: kafka}
}

// Listen starts listening to the "trip_created" topic and processes incoming messages.
func (tec *TripEventConsumer) Listen() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	err := tec.Kafka.SubscribeAndConsume(ctx, "trip_created",
		func(ctx context.Context, msg *kafka.Message) error {
			return nil
		},
	)
	tec.Kafka.Close()
	return err
}
