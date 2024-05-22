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

	group        *gin.RouterGroup
	loginHandler gin.HandlerFunc
	authChecker  auth.AuthorizationChecker
}

func (con *UserController) ConfigureApi(r *gin.RouterGroup) {
	con.group = r.Group("/user")
	con.group.Use(con.loginHandler)
	{
		// TODO
		con.group.GET("/login-test", func(ctx *gin.Context) {
			ctx.IndentedJSON(http.StatusOK, gin.H{
				"message": "hello:)",
			})
		})
	}

	con.authChecker = auth.NewAuthorizationCheckerBuilder().
		ForPath(con.group.BasePath() + "/*").
		ForMethod("*").
		PermitAll().
		Build()
}

func (con UserController) ConfigurePages(r *gin.RouterGroup) {
	// TODO
}

func (con *UserController) Check(c *gin.Context, user *model.User) (authorized bool, matches bool) {
	return con.authChecker.Check(c, user)
}

func NewUserController(userService service.UserService, loginHandler gin.HandlerFunc) *UserController {
	return &UserController{
		userService:  userService,
		loginHandler: loginHandler,
	}
}