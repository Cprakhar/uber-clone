package events

import (
	"time"

	"github.com/cprakhar/uber-clone/shared/messaging"
)

type TripEventProducer struct {
	Kafka *messaging.Producer
}

// NewTripEventProducer creates a new TripEventProducer with the given Kafka producer.
func NewTripEventProducer(kafka *messaging.Producer) *TripEventProducer {
	return &TripEventProducer{Kafka: kafka}
}

// PublishTripCreatedEventAndWait publishes and waits for broker delivery (acks) up to timeout.
func (tep *TripEventProducer) PublishTripCreatedEventAndWait(payload []byte, entityID string, timeout time.Duration) error {
	return tep.Kafka.SendMessageAndWait("trip_created", payload, entityID, timeout)
}
