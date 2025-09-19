package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/cprakhar/uber-clone/services/trip-service/events"
	"github.com/cprakhar/uber-clone/services/trip-service/handler"
	"github.com/cprakhar/uber-clone/services/trip-service/service"
	"github.com/cprakhar/uber-clone/shared/messaging/kafka"
	"google.golang.org/grpc"
)

type gRPCServer struct {
	addr        string
	tripService service.TripService
	kfClient    *kafka.KafkaClient
}

// NewgRPCServer creates a new gRPC server instance
func NewgRPCServer(addr string, tripService service.TripService, kfc *kafka.KafkaClient) *gRPCServer {
	return &gRPCServer{addr: addr, tripService: tripService, kfClient: kfc}
}

// run starts the gRPC server and listens for incoming requests
func (s *gRPCServer) run(ctx context.Context) error {
	// Start listening on the specified address
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %v", s.addr, err)
	}
	
	// gRPC server setup
	srv := grpc.NewServer()
	handler.NewgRPCHandler(srv, s.tripService, events.NewTripEventProducer(s.kfClient))

	// Graceful shutdown on context cancellation
	go func() {
		<-ctx.Done()
		srv.GracefulStop()
	}()

	// Start serving
	log.Printf("gRPC server running on %s", s.addr)
	if err := srv.Serve(lis); err != nil && ctx.Err() == nil {
		return fmt.Errorf("failed to serve gRPC server: %v", err)
	}
	return nil
}
