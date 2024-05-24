package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"store.api/auth"
	"store.api/model"
	"store.api/service"
)

type UserController struct {
	userService service.UserService

	group       *gin.RouterGroup
	auth        gin.HandlerFunc
	authChecker auth.AuthorizationChecker
}

func (con *UserController) ConfigureApi(r *gin.RouterGroup) {
	con.group = r.Group("/user")
	con.group.Use(con.auth)
	{
		con.group.GET("/login-test", func(ctx *gin.Context) {
			ctx.IndentedJSON(http.StatusOK, gin.H{
				"message": "hello:)",
			})
		})
		// TODO
	}

	con.authChecker = auth.NewAuthorizationCheckerBuilder().
		ForPath(con.group.BasePath() + "*").
		ForMethod("*").
		PermitAll().
		Build()
}

func (con *UserController) Check(c *gin.Context, user *model.User) (authorized bool, matches bool) {
	return con.authChecker.Check(c, user)
}

func NewUserController(userService service.UserService, auth gin.HandlerFunc) *UserController {
	return &UserController{
		userService: userService,
		auth:        auth,
	}
}
