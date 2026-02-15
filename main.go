package main

import (
	"github.com/lirajoaop/gopportunities/config"
	"github.com/lirajoaop/gopportunities/router"
)

var (
	logger *config.Logger
)

// @title Go Opportunities API
// @version 1.0
// @description API to manage job opportunities, including creating, listing, updating, and deleting job openings.
// @contact.name João Pedro Lira
// @contact.url https://github.com/lirajoaop/gopportunities
// @BasePath /api/v1
// @schemes http https
func main() {
	logger = config.GetLogger("main")

	err := config.Init()
	if err != nil {
		logger.Errorf("Config initialization error: %v", err)
		return
	}

	router.Initialize()
}
