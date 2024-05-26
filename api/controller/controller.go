package controller

import "github.com/gin-gonic/gin"

type Controller interface {
	ConfigureApi(*gin.RouterGroup)
}
