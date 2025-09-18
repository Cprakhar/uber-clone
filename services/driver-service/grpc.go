package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/cprakhar/uber-clone/services/driver-service/handler"
	"github.com/cprakhar/uber-clone/services/driver-service/repo"
	"github.com/cprakhar/uber-clone/services/driver-service/service"
	"github.com/cprakhar/uber-clone/shared/events"
	"github.com/cprakhar/uber-clone/shared/messaging"
	"google.golang.org/grpc"
)

type gRPCServer struct {
	addr string
}

func NewgRPCServer(addr string) *gRPCServer {
	return &gRPCServer{addr: addr}
}

func (s *gRPCServer) run() {
	driverRepo := repo.NewDriverRepository()
	driverService := service.NewgRPCService(driverRepo)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit
		cancel()
	}()
	
	// Start listening on the specified address
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", s.addr, err)
	}
	// Implement gRPC server setup and start logic here
	srv := grpc.NewServer()
	handler.NewgRPCHandler(srv, driverService)

	// Kafka consumer connection
	consumer, err := messaging.NewConsumer("kafka:9092", "driver-service-group")
	if err != nil {
		log.Fatalf("failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()
	log.Println("Kafka consumer connected")

	// Subscribe to topics and start consuming messages
	tripEventConsumer := events.NewTripEventConsumer(consumer)

	// Start listening to trip events
	go func() {
		if err := tripEventConsumer.Listen(); err != nil {
			log.Fatalf("failed to start trip event consumer: %v", err)
		}
	}()

	// Start gRPC server
	go func() {
		log.Printf("gRPC server started listening on %s", s.addr)
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC server: %v", err)
		}
	}()
	<-ctx.Done()
	log.Println("shutting down gRPC server...")
	srv.GracefulStop()
	log.Println("gRPC server stopped")
}