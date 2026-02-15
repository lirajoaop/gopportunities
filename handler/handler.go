package handler

import (
	"github.com/lirajoaop/gopportunities/config"
	"github.com/lirajoaop/gopportunities/repository"
	"github.com/lirajoaop/gopportunities/service"
)

var (
	logger         *config.Logger
	openingService *service.OpeningService
)

func InitializeHandler() {
	logger = config.GetLogger("handler")
	db := config.GetDatabase()
	openingRepo := repository.NewOpeningRepository(db)
	openingService = service.NewOpeningService(openingRepo)
}
