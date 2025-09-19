package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cprakhar/uber-clone/services/driver-service/events"
	"github.com/cprakhar/uber-clone/services/driver-service/repo"
	"github.com/cprakhar/uber-clone/services/driver-service/service"
	"github.com/cprakhar/uber-clone/shared/contracts"
	"github.com/cprakhar/uber-clone/shared/messaging/kafka"
)

var (
	brokers = []string{"kafka:9092"}
	groupID = "driver-service-group"
	topics  = []string{contracts.TripEventCreated, contracts.TripEventDriverNotInterested}
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize Kafka client
	kfClient, err := kafka.NewKafkaClient(brokers, groupID)
	if err != nil {
		log.Fatalf("failed to create Kafka client: %v", err)
	}
	defer kfClient.Close()
	log.Println("Kafka client connected")

	// Initialize repositories and services
	driverRepo := repo.NewDriverRepository()
	driverService := service.NewDriverService(driverRepo)

	// Start consuming trip events
	tripConsumer := events.NewTripConsumer(kfClient, driverService)
	go func() {
		if err := tripConsumer.Consume(ctx, topics); err != nil {
			log.Printf("Error consuming trip topics: %v", err)
		}
		<-ctx.Done()
	}()

	// Start gRPC server
	grpcServer := NewgRPCServer(":9100", kfClient, driverService)
	go func() {
		if err := grpcServer.run(ctx); err != nil && ctx.Err() == nil {
			log.Printf("gRPC server error: %v", err)
			stop()
		}
	}()

	<-ctx.Done()
	log.Println("Shutdown signal received, exiting...")
}
