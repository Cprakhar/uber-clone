package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cprakhar/uber-clone/shared/contracts"
	"github.com/cprakhar/uber-clone/shared/env"
	"github.com/cprakhar/uber-clone/shared/messaging"
	"github.com/cprakhar/uber-clone/shared/messaging/kafka"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8080")
	brokers  = []string{"kafka:9092"}
	groupID  = "api-gateway-group"
	topics   = []string{
		contracts.TripEventNoDriversFound,
		contracts.TripEventDriverAssigned,
		contracts.DriverCmdTripRequest,
		contracts.PaymentEventSessionCreated,
	}

	connManager = messaging.NewConnectionManager()
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

	topicConsumer := messaging.NewTopicConsumer(kfClient, connManager, topics)
	go func() {
		if err := topicConsumer.Consume(ctx); err != nil && ctx.Err() == nil {
			log.Printf("Error consuming topics: %v", err)
		}
		stop()
	}()

	// Start http server
	httpServer := NewhttpServer(httpAddr, kfClient, connManager)
	go func() {
		if err := httpServer.run(ctx); err != nil && ctx.Err() == nil {
			log.Printf("http server error: %v", err)
			stop()
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	log.Println("Shutdown signal received, exiting...")
}
