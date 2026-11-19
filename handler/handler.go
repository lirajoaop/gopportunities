package handler

import (
	"github.com/lirajoaop/gopportunities/config"
	"github.com/lirajoaop/gopportunities/repository"
	"github.com/lirajoaop/gopportunities/schemas"
	"github.com/lirajoaop/gopportunities/service"
)

// OpeningService is the business-logic contract the handlers depend on.
// It is satisfied by *service.OpeningService and by test fakes.
type OpeningService interface {
	CreateOpening(opening *schemas.Opening) error
	GetOpeningByID(id string) (*schemas.Opening, error)
	ListOpenings() ([]schemas.Opening, error)
	UpdateOpening(id string, role, company, location, link string, remote *bool, salary int64) (*schemas.Opening, error)
	DeleteOpening(id string) (*schemas.Opening, error)
}

var (
	logger         *config.Logger
	openingService OpeningService
)

func InitializeHandler() {
	logger = config.GetLogger("handler")
	db := config.GetDatabase()
	openingRepo := repository.NewOpeningRepository(db)
	openingService = service.NewOpeningService(openingRepo)
}
