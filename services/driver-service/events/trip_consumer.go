package events

import (
	"context"
	"encoding/json"
	"log"
	"math/rand/v2"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/cprakhar/uber-clone/services/driver-service/service"
	"github.com/cprakhar/uber-clone/shared/contracts"
	"github.com/cprakhar/uber-clone/shared/messaging"
	kf "github.com/cprakhar/uber-clone/shared/messaging/kafka"
)

type TripConsumer struct {
	kfClient *kf.KafkaClient
	svc      service.DriverService
}

// NewTripEventConsumer creates a new TripEventConsumer with the given Kafka consumer.
func NewTripConsumer(kfClient *kf.KafkaClient, svc service.DriverService) *TripConsumer {
	return &TripConsumer{kfClient: kfClient, svc: svc}
}

// Consume starts consuming to the specified topics and processes messages.
func (tec *TripConsumer) Consume(ctx context.Context, topics []string) error {
	defer tec.kfClient.Consumer.Close()
	return tec.kfClient.Consumer.SubscribeAndConsume(ctx, topics,
		func(ctx context.Context, msg *kafka.Message) error {

			var kafkaMsg contracts.KafkaMessage
			if err := json.Unmarshal(msg.Value, &kafkaMsg); err != nil {
				log.Printf("failed to unmarshal message: %v", err)
				return err
			}

			log.Printf("Received message on topic %s: %+v", *msg.TopicPartition.Topic, kafkaMsg)

			var payload messaging.TripEventData
			if err := json.Unmarshal(kafkaMsg.Data, &payload); err != nil {
				log.Printf("failed to unmarshal payload: %v", err)
			}

			// Handle different event types
			switch *msg.TopicPartition.Topic {
			case contracts.TripEventCreated, contracts.TripEventDriverNotInterested:
				return tec.handleFindAndNotifyDrivers(ctx, &payload)
			}

			log.Printf("Unknown trip event: %v", payload)
			return nil
		},
	)
}

func (tec *TripConsumer) handleFindAndNotifyDrivers(ctx context.Context, payload *messaging.TripEventData) error {
	drivers := tec.svc.FindAvailableDrivers(ctx, payload.Trip.SelectedFare.PackageSlug)
	log.Printf("Available drivers for trip %s: %v", payload.Trip.Id, drivers)
	if len(drivers) == 0 {
		log.Printf("No drivers available for trip %s", payload.Trip.Id)

		// Notify trip service about unavailability of drivers
		if err := tec.kfClient.Producer.SendMessage(contracts.TripEventNoDriversFound, &contracts.KafkaMessage{
			EntityID: payload.Trip.Id,
		}); err != nil {
			log.Printf("failed to notify trip service about no drivers found: %v", err)
		}
		return nil
	}

	randIdx := rand.IntN(len(drivers))
	selectedDriverID := drivers[randIdx]

	marshalledEvent, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal data: %v", err)
	}

	// Notify trip service about the selected driver
	if err := tec.kfClient.Producer.SendMessage(contracts.DriverCmdTripRequest, &contracts.KafkaMessage{
		EntityID: selectedDriverID,
		Data:     marshalledEvent,
	}); err != nil {
		log.Printf("failed to notify trip service about selected driver: %v", err)
		return err
	}
	log.Printf("Found a suitable driver %s for trip %s", selectedDriverID, payload.Trip.Id)
	return nil
}
