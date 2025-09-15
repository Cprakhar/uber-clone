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

	srv := grpc.NewServer()
	handler.NewgRPCHandler(srv, tripService)

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down gRPC server...")
	srv.GracefulStop()
}
