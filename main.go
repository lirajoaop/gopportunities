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
// @description API to manage job opportunities
// @host localhost:8080
// @BasePath /api/v1
func main() {
	logger = config.GetLogger("main")

	err := config.Init()
	if err != nil {
		logger.Errorf("Config initialization error: %v", err)
		return
	}

	router.Initialize()
}
