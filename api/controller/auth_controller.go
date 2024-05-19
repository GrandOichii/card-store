package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"store.api/dto"
	"store.api/service"
)

type AuthController struct {
	Controller

	userService service.UserService
}

func (con AuthController) Configure(r *gin.RouterGroup) {
	// g := r.Group("/api/v1/auth")
	// {
	r.POST("/auth/register", con.Register)
	// g.POST("/login", con.Login)
	// }
}

func NewAuthController(userService service.UserService) *AuthController {
	return &AuthController{
		userService: userService,
	}
}

// UserRegister			godoc
// @Summary				Registers the user
// @Description			Checks the user data and adds it to the repo
// @Param				details body dto.RegisterDetails true "Register details"
// @Tags				Auth
// @Success				200 {object} dto.PrivateUserInfo
// @Router				/auth/register [post]
func (con *AuthController) Register(c *gin.Context) {
	var newUser dto.RegisterDetails

	if err := c.BindJSON(&newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := con.userService.Register(&newUser)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, user)
}
