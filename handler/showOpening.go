package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Show opening
// @Description Get a job opening by ID
// @Tags Opening
// @Accept json
// @Produce json
// @Param id path string true "Opening ID"
// @Success 200 {object} ShowOpeningResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /openings/{id} [get]
func ShowOpeningHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		sendError(ctx, http.StatusBadRequest, errParamIsRequired("id", "pathParameter").Error())
		return
	}

	opening, err := openingService.GetOpeningByID(id)
	if err != nil {
		sendError(ctx, http.StatusNotFound, err.Error())
		return
	}

	sendSuccess(ctx, "show-opening", opening.ToResponse())
}
