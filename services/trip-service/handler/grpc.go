package handler

import (
	"context"
	"log"

	"github.com/cprakhar/uber-clone/services/trip-service/events"
	"github.com/cprakhar/uber-clone/services/trip-service/service"
	"github.com/cprakhar/uber-clone/services/trip-service/types"
	pb "github.com/cprakhar/uber-clone/shared/proto/trip"
	sharedtypes "github.com/cprakhar/uber-clone/shared/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gRPCHandler struct {
	pb.UnimplementedTripServiceServer
	svc      service.TripService
	producer *events.TripEventProducer
}

// NewgRPCHandler registers the gRPC handler with the given gRPC server
func NewgRPCHandler(srv *grpc.Server, svc service.TripService, producer *events.TripEventProducer) {
	handler := &gRPCHandler{svc: svc, producer: producer}
	pb.RegisterTripServiceServer(srv, handler)
}

// PreviewTrip handles the PreviewTrip gRPC request
func (h *gRPCHandler) PreviewTrip(ctx context.Context, req *pb.PreviewTripRequest) (*pb.PreviewTripResponse, error) {
	pickup := req.GetPickup()
	destination := req.GetDestination()

	pickupCoords := &sharedtypes.Coordinate{
		Latitude:  pickup.GetLatitude(),
		Longitude: pickup.GetLongitude(),
	}
	destinationCoords := &sharedtypes.Coordinate{
		Latitude:  destination.GetLatitude(),
		Longitude: destination.GetLongitude(),
	}

	route, err := h.svc.GetRoute(ctx, pickupCoords, destinationCoords)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get route: %v", err)
	}

	estimatedFares := h.svc.EstimatePackagesPriceWithRoute(route)

	fares, err := h.svc.GenerateTripFares(ctx, estimatedFares, req.GetRiderID(), route)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate trip fares: %v", err)
	}

	return &pb.PreviewTripResponse{
		Route:     route.ToProto(),
		RideFares: types.ToRideFaresProto(fares),
	}, nil
}

// CreateTrip handles the CreateTrip gRPC request
func (h *gRPCHandler) CreateTrip(ctx context.Context, req *pb.CreateTripRequest) (*pb.CreateTripResponse, error) {
	fareID := req.GetRideFareID()
	riderID := req.GetRiderID()
	fare, err := h.svc.GetAndValidateRideFare(ctx, fareID, riderID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid fare: %v", err)
	}

	trip, err := h.svc.CreateTrip(ctx, fare)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create trip: %v", err)
	}

	// Notify other services about the new trip
	if err := h.producer.PublishTripCreated(trip); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to publish trip created event: %v", err)
	}
	log.Printf("Published trip created event for trip ID: %s", trip.ID.Hex())

	return &pb.CreateTripResponse{
		TripID: trip.ID.Hex(),
	}, nil
}
