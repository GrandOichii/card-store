package main

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"store.api/config"
	docs "store.api/docs"
	"store.api/router"
)

// @title Card store api
// @version 1.0
// @description Service for buying/selling collectable cards
// @host localhost:8080
// @BasePath /api/v1
func main() {
	var c *config.Configuration

	c, err := config.ReadConfig("config.json")

	if err != nil {
		c, err = config.ReadEnvConfig()
		if err != nil {
			panic(err)
		}
	}

	router := router.CreateRouter(c)

	docs.SwaggerInfo.BasePath = "/api/v1"
	if gin.Mode() != gin.ReleaseMode {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
	router.Run(c.Host + ":" + c.Port)
}
