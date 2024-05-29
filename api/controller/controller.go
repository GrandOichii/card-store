package controller

import (
	"github.com/gin-gonic/gin"
)

type Controller interface {
	ConfigureApi(*gin.RouterGroup)
}

func AbortWithError(c *gin.Context, status int, err error, public bool) {
	if public {
		c.AbortWithStatusJSON(status, err.Error())
		return
	}
	c.AbortWithError(status, err)
}
