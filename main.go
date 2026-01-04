package main

import (
	"github.com/lirajoaop/gopportunities/config"
	"github.com/lirajoaop/gopportunities/router"
)

var (
	logger *config.Logger
)

func main() {
	logger = config.GetLogger("main")

	err := config.Init()
	if err != nil {
		logger.Errorf("Config initialization error: %v", err)
		return
	}

	router.Initialize()
}
