package types

import (
	pb "github.com/cprakhar/uber-clone/shared/proto/trip"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripModel struct {
	ID       primitive.ObjectID
	RiderID  string
	Status   string
	RideFare *RideFareModel
	Driver   *pb.TripDriver
}

type RideFareModel struct {
	ID               primitive.ObjectID
	RiderID          string
	PackageSlug      string
	TotalFareInPaise float64
	Route            *OSRMApiResponse
}

func (r *RideFareModel) ToProto() *pb.RideFare {
	return &pb.RideFare{
		Id:               r.ID.Hex(),
		RiderID:          r.RiderID,
		PackageSlug:      r.PackageSlug,
		TotalFareInPaise: r.TotalFareInPaise,
	}
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
			Latitude:  coord[0],
			Longitude: coord[1],
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

func ToRideFaresProto(fares []*RideFareModel) []*pb.RideFare {
	protoFares := make([]*pb.RideFare, len(fares))
	for i, fare := range fares {
		protoFares[i] = fare.ToProto()
	}
	return protoFares
}

type PricingConfig struct {
	PricePerUnitDistance float64
	PricePerMinute       float64
}

func DefaultPricingConfig() *PricingConfig {
	return &PricingConfig{
		PricePerUnitDistance: 10.0, // 10 Rs per km
		PricePerMinute:       2.0,  // 2 Rs per minute
	}
}
