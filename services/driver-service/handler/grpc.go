package handler

import (
	"context"
	"log"

	"github.com/cprakhar/uber-clone/services/driver-service/service"
	pb "github.com/cprakhar/uber-clone/shared/proto/driver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gRPCHandler struct {
	pb.UnimplementedDriverServiceServer
	svc service.DriverService
}

func NewgRPCHandler(srv *grpc.Server, svc service.DriverService) {
	handler := &gRPCHandler{svc: svc}
	pb.RegisterDriverServiceServer(srv, handler)
}

func (h *gRPCHandler) RegisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	driverID := req.GetDriverID()
	packageSlug := req.GetPackageSlug()
	// Implement the logic to register a driver
	driver, err := h.svc.RegisterDriver(ctx, driverID, packageSlug)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register driver: %v", err)
	}
	log.Printf("Driver registered: %s", driver.Id)
	return &pb.RegisterDriverResponse{
		Driver: driver,
	}, nil
}

func (h *gRPCHandler) UnregisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	// Implement the logic to unregister a driver
	return nil, nil
}
