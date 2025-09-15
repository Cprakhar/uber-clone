package handler

import (
	"context"

	"github.com/cprakhar/uber-clone/services/trip-service/service"
	pb "github.com/cprakhar/uber-clone/shared/proto/trip"
	"github.com/cprakhar/uber-clone/shared/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gRPCHandler struct {
	pb.UnimplementedTripServiceServer
	svc service.GRPCTripService
}

func NewgRPCHandler(srv *grpc.Server, svc service.GRPCTripService) *gRPCHandler {
	handler := &gRPCHandler{svc: svc}
	pb.RegisterTripServiceServer(srv, handler)
	return handler
}

func (h *gRPCHandler) PreviewTrip(ctx context.Context, req *pb.PreviewTripRequest) (*pb.PreviewTripResponse, error) {
	pickup := req.GetPickup()
	destination := req.GetDestination()

	pickupCoords := &types.Coordinate{
		Latitude:  pickup.GetLatitude(),
		Longitude: pickup.GetLongitude(),
	}
	destinationCoords := &types.Coordinate{
		Latitude:  destination.GetLatitude(),
		Longitude: destination.GetLongitude(),
	}

	route, err := h.svc.GetRoute(ctx, pickupCoords, destinationCoords)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get route: %v", err)
	}

	return &pb.PreviewTripResponse{
		Route:    route.ToProto(),
		RideFare: []*pb.RideFare{},
	}, nil
}
