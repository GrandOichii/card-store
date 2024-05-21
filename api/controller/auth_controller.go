package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"store.api/dto"
	"store.api/service"
)

type AuthController struct {
	userService  service.UserService
	loginHandler gin.HandlerFunc
	group        *gin.RouterGroup
}

func (con AuthController) ConfigureApi(r *gin.RouterGroup) {
	con.group = r.Group("/auth")
	{
		con.group.POST("/register", con.Register)
		con.group.POST("/login", con.Login)
	}
}

func (con AuthController) ConfigureViews(r *gin.RouterGroup) {
	// TODO
}

func NewAuthController(userService service.UserService, loginHandler gin.HandlerFunc) *AuthController {
	return &AuthController{
		userService:  userService,
		loginHandler: loginHandler,
	}
}

// UserRegister			godoc
// @Summary				Registers the user
// @Description			Checks the user data and adds it to the repo
// @Param				details body dto.RegisterDetails true "Register details"
// @Tags				Auth
// @Success				201
// @Router				/auth/register [post]
func (con *AuthController) Register(c *gin.Context) {
	var newUser dto.RegisterDetails

	if err := c.BindJSON(&newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := con.userService.Register(&newUser)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusCreated)
}

// UserLogin			godoc
// @Summary				Logs in the user
// @Description			Checks the user data and returns a jwt token on correct Login
// @Param				details body dto.LoginDetails true "Login details"
// @Tags				Auth
// @Success				200
// @Router				/auth/login [post]
func (con AuthController) Login(c *gin.Context) {
	con.loginHandler(c)
}
