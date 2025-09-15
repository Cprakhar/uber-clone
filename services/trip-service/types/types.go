package types

import (
	"time"

	pb "github.com/cprakhar/uber-clone/shared/proto/trip"
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

type OSRMApiResponse struct {
	Routes []struct {
		Distance float64 `json:"distance"`
		Duration float64 `json:"duration"`
		Geometry struct {
			Coordinates [][]float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"routes"`
}

func (o *OSRMApiResponse) ToProto() *pb.Route {
	if len(o.Routes) == 0 {
		return &pb.Route{}
	}

	route := o.Routes[0]
	coordinates := make([]*pb.Coordinate, len(route.Geometry.Coordinates))
	for i, coord := range route.Geometry.Coordinates {
		if len(coord) != 2 {
			continue
		}
		coordinates[i] = &pb.Coordinate{
			Longitude: coord[0],
			Latitude:  coord[1],
		}
	}
	return &pb.Route{
		Distance: route.Distance,
		Duration: route.Duration,
		Geometry: []*pb.Geometry{
			{
				Coordinates: coordinates,
			},
		},
	}
}
