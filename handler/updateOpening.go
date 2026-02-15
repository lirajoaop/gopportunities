package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Update opening
// @Description Update a job opening by ID
// @Tags Opening
// @Accept json
// @Produce json
// @Param id path string true "Opening ID"
// @Param request body UpdateOpeningRequest true "Request Body"
// @Success 200 {object} UpdateOpeningResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /openings/{id} [put]
func UpdateOpeningHandler(ctx *gin.Context) {
	request := UpdateOpeningRequest{}

	ctx.BindJSON(&request)

	if err := request.Validate(); err != nil {
		logger.Errorf("validation error: %v", err.Error())
		sendError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id := ctx.Param("id")

	if id == "" {
		sendError(ctx, http.StatusBadRequest, errParamIsRequired("id", "pathParameter").Error())
		return
	}

	opening, err := openingService.UpdateOpening(
		id, request.Role, request.Company, request.Location,
		request.Link, request.Remote, request.Salary,
	)
	if err != nil {
		logger.Errorf("error updating opening: %v", err.Error())
		sendError(ctx, http.StatusNotFound, err.Error())
		return
	}

	sendSuccess(ctx, "update-opening", opening.ToResponse())
}
