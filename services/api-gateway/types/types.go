package types

import (
	pb "github.com/cprakhar/uber-clone/shared/proto/trip"
	"github.com/cprakhar/uber-clone/shared/types"
)

type PreviewTripRequest struct {
	RiderID     string           `json:"riderID" binding:"required"`
	Pickup      types.Coordinate `json:"pickup" binding:"required"`
	Destination types.Coordinate `json:"destination" binding:"required"`
}

// ToProto converts PreviewTripRequest to its protobuf representation
func (ptr *PreviewTripRequest) ToProto() *pb.PreviewTripRequest {
	return &pb.PreviewTripRequest{
		RiderID: ptr.RiderID,
		Pickup: &pb.Coordinate{
			Latitude:  ptr.Pickup.Latitude,
			Longitude: ptr.Pickup.Longitude,
		},
		Destination: &pb.Coordinate{
			Latitude:  ptr.Destination.Latitude,
			Longitude: ptr.Destination.Longitude,
		},
	}
}

type TripStartRequest struct {
	RiderID string `json:"riderID" binding:"required"`
	FareID  string `json:"rideFareID" binding:"required"`
}

// ToProto converts TripStartRequest to its protobuf representation
func (tsr *TripStartRequest) ToProto() *pb.CreateTripRequest {
	return &pb.CreateTripRequest{
		RiderID:    tsr.RiderID,
		RideFareID: tsr.FareID,
	}
}
