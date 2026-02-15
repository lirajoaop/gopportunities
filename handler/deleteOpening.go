package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Delete opening
// @Description Delete a job opening by ID
// @Tags Opening
// @Accept json
// @Produce json
// @Param id path string true "Opening ID"
// @Success 200 {object} DeleteOpeningResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /openings/{id} [delete]
func DeleteOpeningHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		sendError(ctx, http.StatusBadRequest, errParamIsRequired("id", "pathParameter").Error())
		return
	}

	opening, err := openingService.DeleteOpening(id)
	if err != nil {
		logger.Errorf("error deleting opening: %v", err.Error())
		sendError(ctx, http.StatusNotFound, err.Error())
		return
	}

	sendSuccess(ctx, "delete-opening", opening.ToResponse())
}
