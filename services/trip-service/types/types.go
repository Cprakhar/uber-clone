package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripModel struct {
	ID       primitive.ObjectID
	RiderID  primitive.ObjectID
	Status   string
	RideFare *RideFareModel
}

type RideFareModel struct {
	ID                primitive.ObjectID
	RiderID           primitive.ObjectID
	PackageSlug       string
	TotalFareInRupees float64
	ExpiresAt         time.Time
}
