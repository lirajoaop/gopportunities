package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lirajoaop/gopportunities/schemas"
)

func sendError(ctx *gin.Context, code int, msg string) {
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(code, gin.H{
		"message":   msg,
		"errorCode": code,
	})
}

func sendSuccess(ctx *gin.Context, op string, data interface{}) {
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(200, gin.H{
		"message": fmt.Sprintf("operation %s successful", op),
		"data":    data,
	})
}

type ErrorResponse struct {
	Message   string `json:"message" example:"error message"`
	ErrorCode int    `json:"errorCode" example:"400"`
}

type CreateOpeningResponse struct {
	Message string                  `json:"message" example:"operation create-opening successful"`
	Data    schemas.OpeningResponse `json:"data"`
}

type DeleteOpeningResponse struct {
	Message string                  `json:"message" example:"operation delete-opening successful"`
	Data    schemas.OpeningResponse `json:"data"`
}

type ShowOpeningResponse struct {
	Message string                  `json:"message" example:"operation show-opening successful"`
	Data    schemas.OpeningResponse `json:"data"`
}

type ListOpeningsResponse struct {
	Message string                    `json:"message" example:"operation list-openings successful"`
	Data    []schemas.OpeningResponse `json:"data"`
}

type UpdateOpeningResponse struct {
	Message string                  `json:"message" example:"operation update-opening successful"`
	Data    schemas.OpeningResponse `json:"data"`
}
