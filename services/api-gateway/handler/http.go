package handler

import (
	"log"
	"net/http"

	grpcclient "github.com/cprakhar/uber-clone/services/api-gateway/grpc-client"
	"github.com/cprakhar/uber-clone/services/api-gateway/types"
	"github.com/cprakhar/uber-clone/shared/contracts"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NewHTTPHandler initializes the HTTP handler with routes and middleware
func NewHTTPHandler() *gin.Engine {
	r := gin.Default()

	// Health check endpoints
	r.GET("/health", healthHandler)
	r.GET("/ready", readinessHandler)

	// Use CORS middleware
	// Define your middleware here
	// Define your routes and handlers here
	api := r.Group("/api/v1")
	{
		api.POST("/trips/preview", enableCORS, previewTripHandler)
		api.GET("/ws/riders", RidersWSHandler)
		api.GET("/ws/drivers", DriversWSHandler)
	}

	return r
}

// healthHandler handles liveness probe requests
func healthHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "service": "api-gateway"})
}

// readinessHandler handles readiness probe requests
func readinessHandler(ctx *gin.Context) {
	// Here you could add checks for dependencies like databases, external services
	ctx.JSON(http.StatusOK, gin.H{"status": "ready", "service": "api-gateway"})
}

// previewTripHandler handles trip preview requests
func previewTripHandler(ctx *gin.Context) {
	var payload types.PreviewTripRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, contracts.APIResponse{
			Error: &contracts.APIError{
				Code:    http.StatusBadRequest,
				Message: "Invalid request payload",
			},
		})
		return
	}

	if payload.RiderID == primitive.NilObjectID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "rider_id is required"})
		return
	}

	tripService, err := grpcclient.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer tripService.Close()

	tripPreview, err := tripService.Client.PreviewTrip(ctx, payload.ToProto())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to preview trip"})
		return
	}
	
	res := contracts.APIResponse{Data: tripPreview}
	ctx.JSON(http.StatusOK, res)
}

func enableCORS(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if ctx.Request.Method == "OPTIONS" {
		ctx.Writer.WriteHeader(http.StatusOK)
		return
	}

	ctx.Next()
}
