package controller

import (
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
	con.group = r.Group("/me")
	con.group.Use(con.loginHandler)
	{
		// TODO
	}

	con.authChecker = auth.NewAuthorizationCheckerBuilder(con.group.BasePath()).
		ForMethod("*").
		PermitAll().
		Build()
}

func (con UserController) ConfigureViews(r *gin.RouterGroup) {
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
