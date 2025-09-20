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

type PaymentEventSessionCreatedData struct {
	TripID    string  `json:"tripID"`
	SessionID string  `json:"sessionID"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
}

type PaymentTripResponseData struct {
	TripID   string  `json:"tripID"`
	RiderID  string  `json:"riderID"`
	DriverID string  `json:"driverID"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type PaymentStatusUpdateData struct {
	TripID   string `json:"tripID"`
	RiderID  string `json:"riderID"`
	DriverID string `json:"driverID"`
}
