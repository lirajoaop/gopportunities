package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lirajoaop/gopportunities/schemas"
)

// @Summary List all openings
// @Description Get all job openings
// @Tags Opening
// @Accept json
// @Produce json
// @Success 200 {object} ListOpeningsResponse
// @Failure 500 {object} ErrorResponse
// @Router /openings [get]
func ListOpeningsHandler(ctx *gin.Context) {
	openings, err := openingService.ListOpenings()
	if err != nil {
		sendError(ctx, http.StatusInternalServerError, "error listing openings")
		return
	}

	response := make([]schemas.OpeningResponse, len(openings))
	for i, opening := range openings {
		response[i] = opening.ToResponse()
	}

	sendSuccess(ctx, "list-openings", response)
}
