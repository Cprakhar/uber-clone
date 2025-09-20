package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cprakhar/uber-clone/services/payment-service/events"
	"github.com/cprakhar/uber-clone/services/payment-service/service"
	"github.com/cprakhar/uber-clone/services/payment-service/types"
	"github.com/cprakhar/uber-clone/shared/contracts"
	"github.com/cprakhar/uber-clone/shared/env"
	"github.com/cprakhar/uber-clone/shared/messaging/kafka"
)

var (
	brokers = []string{"kafka:9092"}
	groupID = "payment-service-group"
	appURL  = env.GetString("APP_URL", "http://localhost:3000")
	topics  = []string{contracts.PaymentCmdCreateSession}
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	kfClient, err := kafka.NewKafkaClient(brokers, groupID)
	if err != nil {
		log.Fatalf("Failed to create Kafka client: %v", err)
	}
	defer kfClient.Close()
	log.Println("Kafka client connected")

	stripeCfg := &types.PaymentConfig{
		StripeSecretKey: env.GetString("STRIPE_SECRET_KEY", ""),
		SuccessURL:      env.GetString("STRIPE_SUCCESS_URL", appURL+"?payment=success"),
		CancelURL:       env.GetString("STRIPE_CANCEL_URL", appURL+"?payment=cancel"),
	}

	if stripeCfg.StripeSecretKey == "" {
		log.Fatal("STRIPE_SECRET_KEY is not set")
	}

	paymentProcessor := service.NewStripeClient(stripeCfg)
	paymentService := service.NewPaymentService(paymentProcessor)

	tripConsumer := events.NewTripConsumer(kfClient, paymentService)
	go func() {
		if err := tripConsumer.Consume(ctx, topics); err != nil && ctx.Err() == nil {
			log.Printf("Error consuming payment topics: %v", err)
		}
		stop()
	}()

	<-ctx.Done()
	log.Println("Shutdown signal received, exiting...")
}
