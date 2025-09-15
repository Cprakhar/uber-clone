package handler

import (
	"net/http"

	"github.com/cprakhar/uber-clone/services/api-gateway/types"
	"github.com/cprakhar/uber-clone/shared/contracts"
	"github.com/gin-gonic/gin"
)

// NewHTTPHandler initializes the HTTP handler with routes and middleware
func NewHTTPHandler() *gin.Engine {
	r := gin.Default()

	// Use CORS middleware
	// Define your middleware here
	// Define your routes and handlers here
	api := r.Group("/api/v1")
	{
		api.POST("/trips/preview", previewTripHandler)
	}

	return r
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

	res := contracts.APIResponse{Data: "Trip previewed successfully"}
	ctx.JSON(http.StatusOK, res)
}
