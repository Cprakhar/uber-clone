package events

import (
	"context"
	"encoding/json"
	"log"

	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/cprakhar/uber-clone/services/trip-service/service"
	"github.com/cprakhar/uber-clone/shared/contracts"
	"github.com/cprakhar/uber-clone/shared/messaging"
	"github.com/cprakhar/uber-clone/shared/messaging/kafka"
	pbd "github.com/cprakhar/uber-clone/shared/proto/driver"
	pb "github.com/cprakhar/uber-clone/shared/proto/trip"
)

type DriverConsumer struct {
	kfClient *kafka.KafkaClient
	svc      service.TripService
}

// NewDriverConsumer creates a new DriverConsumer with the given Kafka consumer.
func NewDriverConsumer(kfClient *kafka.KafkaClient, svc service.TripService) *DriverConsumer {
	return &DriverConsumer{kfClient: kfClient, svc: svc}
}

// Consume starts consuming messages from the specified topics and processes them.
func (dc *DriverConsumer) Consume(ctx context.Context, topics []string) error {
	defer dc.kfClient.Consumer.Close()
	return dc.kfClient.Consumer.SubscribeAndConsume(ctx, topics, func(ctx context.Context, msg *ckafka.Message) error {

		var kafkaMsg contracts.KafkaMessage
		if err := json.Unmarshal(msg.Value, &kafkaMsg); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return err
		}

		var payload messaging.DriverTripResponseData
		if kafkaMsg.Data != nil {
			if err := json.Unmarshal(kafkaMsg.Data, &payload); err != nil {
				log.Printf("Failed to unmarshal payload: %v", err)
				return err
			}
		}

		switch *msg.TopicPartition.Topic {
		case contracts.DriverCmdTripAccept:
			if err := dc.handleTripAccept(ctx, payload.TripID, payload.Driver); err != nil {
				log.Printf("Failed to handle trip accept: %v", err)
				return err
			}
		case contracts.DriverCmdTripDecline:
			log.Printf("Driver declined trip %s", payload.TripID)
			dc.handleTripDecline(ctx, payload.TripID)
		}

		return nil
	})
}

func (dc *DriverConsumer) handleTripDecline(ctx context.Context, tripID string) error {
	trip, err := dc.svc.GetTripByID(ctx, tripID)
	if err != nil {
		log.Printf("Failed to get trip by ID: %v", err)
		return err
	}

	tripEventData := &messaging.TripEventData{
		Trip: trip.ToProto(),
	}

	data, err := json.Marshal(tripEventData)
	if err != nil {
		log.Printf("Failed to marshal trip event data: %v", err)
	}

	// Notify driver service to find another driver
	if err := dc.kfClient.Producer.SendMessage(contracts.TripEventDriverNotInterested, &contracts.KafkaMessage{
		EntityID: trip.RiderID,
		Data:     data,
	}); err != nil {
		log.Printf("Failed to send driver not interested message: %v", err)
	}

	return nil
}

// handleTripAccept processes a trip acceptance from a driver.
func (dc *DriverConsumer) handleTripAccept(ctx context.Context, tripID string, driver *pbd.Driver) error {
	updatedTrip, err := dc.svc.AcceptRide(ctx, tripID, &pb.TripDriver{
		Id:         driver.Id,
		Name:       driver.Name,
		ProfilePic: driver.ProfilePic,
		CarPlate:   driver.CarPlate,
	})
	if err != nil {
		log.Printf("Failed to update trip with driver: %v", err)
		return err
	}

	data, err := json.Marshal(updatedTrip)
	if err != nil {
		log.Printf("Failed to marshal updated trip: %v", err)
		return err
	}

	// Notify rider about driver assignment
	if err := dc.kfClient.Producer.SendMessage(contracts.TripEventDriverAssigned, &contracts.KafkaMessage{
		EntityID: updatedTrip.RiderID,
		Data:     data,
	}); err != nil {
		log.Printf("Failed to send driver assigned message: %v", err)
		return err
	}

	// Notify payment service to initiate payment
	return nil
}
