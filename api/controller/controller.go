package controller

import "github.com/gin-gonic/gin"

type ErrResponse struct {
	Error string `json:"error"`
}

type Controller interface {
	ConfigureApi(*gin.RouterGroup)
}
