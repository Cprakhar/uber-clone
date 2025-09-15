package types

import (
	"github.com/cprakhar/uber-clone/shared/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PreviewTripRequest struct {
	RiderID     primitive.ObjectID `json:"rider_id" binding:"required"`
	Pickup      types.Coordinate   `json:"pickup" binding:"required"`
	Destination types.Coordinate   `json:"destination" binding:"required"`
}
