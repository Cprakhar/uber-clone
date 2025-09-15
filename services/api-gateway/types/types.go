package types

import (
	pb "github.com/cprakhar/uber-clone/shared/proto/trip"
	"github.com/cprakhar/uber-clone/shared/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PreviewTripRequest struct {
	RiderID     primitive.ObjectID `json:"rider_id" binding:"required"`
	Pickup      types.Coordinate   `json:"pickup" binding:"required"`
	Destination types.Coordinate   `json:"destination" binding:"required"`
}

func (ptr *PreviewTripRequest) ToProto() *pb.PreviewTripRequest {
	return &pb.PreviewTripRequest{
		RiderID: ptr.RiderID.Hex(),
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
