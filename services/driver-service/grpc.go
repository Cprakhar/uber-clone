package main

import (
	"context"
	"log"
	"net"

	"github.com/cprakhar/uber-clone/services/driver-service/handler"
	"github.com/cprakhar/uber-clone/services/driver-service/service"
	"github.com/cprakhar/uber-clone/shared/messaging/kafka"
	"google.golang.org/grpc"
)

type gRPCServer struct {
	addr          string
	kfClient      *kafka.KafkaClient
	driverService service.DriverService
}

func NewgRPCServer(addr string, kfc *kafka.KafkaClient, svc service.DriverService) *gRPCServer {
	return &gRPCServer{addr: addr, kfClient: kfc, driverService: svc}
}

func (s *gRPCServer) run(ctx context.Context) error {
	// Start listening on the specified address
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Printf("Failed to listen on %s: %v", s.addr, err)
		return err
	}

	// gRPC server setup
	srv := grpc.NewServer()
	handler.NewgRPCHandler(srv, s.driverService)

	// Graceful shutdown on context cancellation
	go func() {
		<-ctx.Done()
		srv.GracefulStop()
	}()

	// Start serving
	log.Printf("gRPC server running on %s", s.addr)
	if err := srv.Serve(lis); err != nil && ctx.Err() == nil {
		log.Printf("Failed to serve gRPC server: %v", err)
		return err
	}
	return nil
}
