package messaging

import (
	pbd "github.com/cprakhar/uber-clone/shared/proto/driver"
	pb "github.com/cprakhar/uber-clone/shared/proto/trip"
)

type TripEventData struct {
	Trip *pb.Trip `json:"trip"`
}

type DriverTripResponseData struct {
	Driver  *pbd.Driver `json:"driver"`
	RiderID string      `json:"riderID"`
	TripID  string      `json:"tripID"`
}
