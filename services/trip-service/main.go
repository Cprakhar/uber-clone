package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cprakhar/uber-clone/services/trip-service/events"
	"github.com/cprakhar/uber-clone/services/trip-service/repo"
	"github.com/cprakhar/uber-clone/services/trip-service/service"
	"github.com/cprakhar/uber-clone/shared/contracts"
	"github.com/cprakhar/uber-clone/shared/messaging/kafka"
)

var (
	brokers = []string{"kafka:9092"}
	groupID = "trip-service-group"
	topics  = []string{contracts.DriverCmdTripAccept, contracts.DriverCmdTripDecline}
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize Kafka client
	kfClient, err := kafka.NewKafkaClient(brokers, groupID)
	if err != nil {
		log.Fatalf("Failed to create Kafka client: %v", err)
	}
	defer kfClient.Close()
	log.Println("Kafka client connected")

	// Initialize repositories and services
	tripRepo := repo.NewInMemoRepository()
	tripService := service.NewService(tripRepo)

	// Start consuming driver responses
	driverConsumer := events.NewDriverConsumer(kfClient, tripService)
	go func() {
		if err := driverConsumer.Consume(ctx, topics); err != nil {
			log.Printf("Error consuming driver topics: %v", err)
		}
		<-ctx.Done()
	}()

	// Start gRPC server
	gRPCServer := NewgRPCServer(":9000", tripService, kfClient)
	go func() {
		if err := gRPCServer.run(ctx); err != nil && ctx.Err() == nil {
			log.Printf("gRPC server error: %v", err)
			stop()
		}
	}()

	<-ctx.Done()
	log.Println("Shutdown signal received, exiting...")
}
