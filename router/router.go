package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lirajoaop/gopportunities/config"
)

func Initialize() {
	//Initialize Router
	router := gin.Default()

	//Add CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     config.GetCORSOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	//Initialize Routes
	initializeRoutes(router)

	//Run server
	router.Run(config.GetPort())
}
