package auth

import (
	"github.com/gin-gonic/gin"
	"store.api/model"
)

type AuthorizationChecker interface {
	Check(c *gin.Context, user *model.User) (authorized bool, matches bool)
}
