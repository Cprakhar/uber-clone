package events

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cprakhar/uber-clone/services/trip-service/types"
	"github.com/cprakhar/uber-clone/shared/contracts"
	"github.com/cprakhar/uber-clone/shared/messaging"
	"github.com/cprakhar/uber-clone/shared/messaging/kafka"
)

type TripEventProducer struct {
	k *kafka.KafkaClient
}

// NewTripEventProducer creates a new TripEventProducer with the given Kafka producer.
func NewTripEventProducer(k *kafka.KafkaClient) *TripEventProducer {
	return &TripEventProducer{k: k}
}

// PublishTripCreated publishes a "trip.event.created" event with the given payload and entity ID.
func (tep *TripEventProducer) PublishTripCreated(trip *types.TripModel, timeout ...time.Duration) error {
	msg := messaging.TripEventData{
		Trip: trip.ToProto(),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	return tep.k.Producer.SendMessage(contracts.TripEventCreated, &contracts.KafkaMessage{
		EntityID: trip.RiderID,
		Data:     data,
	})
}
