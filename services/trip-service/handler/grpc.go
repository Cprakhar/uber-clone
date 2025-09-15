package handler

import "github.com/cprakhar/uber-clone/services/trip-service/service"

type gRPCHandler struct{
	svc service.GRPCTripService
}

func NewgRPCHandler(svc service.GRPCTripService) *gRPCHandler {
	return &gRPCHandler{svc: svc}
}

func (h *gRPCHandler) CreateTripHandler() {
	
}