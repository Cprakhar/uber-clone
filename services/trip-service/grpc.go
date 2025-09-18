package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/cprakhar/uber-clone/services/trip-service/handler"
	"github.com/cprakhar/uber-clone/services/trip-service/repo"
	"github.com/cprakhar/uber-clone/services/trip-service/service"
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

	tripRepository := repo.NewInMemoRepository()
	tripService := service.NewgRPCService(tripRepository)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit
		cancel()
	}()
		
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", s.addr, err)
	}
		
	// Kafka producer connection
	producer, err := messaging.NewProducer([]string{"kafka:9092"})
	if err != nil {
		log.Fatalf("failed to create Kafka producer: %v", err)
	}
	defer producer.Close()
	
	log.Println("Kafka producer connected")

	// Initialize TripEventProducer
	tripProducer := events.NewTripEventProducer(producer)
	defer tripProducer.Kafka.Close()

	// gRPC server setup
	srv := grpc.NewServer()
	handler.NewgRPCHandler(srv, tripService, tripProducer)

	go func() {
		log.Printf("gRPC server running on %s", s.addr)
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down gRPC server...")
	srv.GracefulStop()
}
