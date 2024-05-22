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

func (con AuthController) ConfigurePages(r *gin.RouterGroup) {
	r.GET("/auth/register", con.RegisterPage)
	r.GET("/auth/login", con.LoginPage)

	// don't like the name
	views := r.Group("/view/auth")
	{
		views.POST("/register", con.RegisterView)
		views.POST("/login", con.LoginView)
	}
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
// @Failure				400
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

	// res, err := con.userService.Login(newUser.ToLoginDetails())
	c.Status(http.StatusCreated)
}

// TODO? more failure docs?
// UserLogin			godoc
// @Summary				Logs in the user
// @Description			Checks the user data and returns a jwt token on correct Login
// @Param				details body dto.LoginDetails true "Login details"
// @Tags				Auth
// @Success				200
// @Failure				400
// @Router				/auth/login [post]
func (con AuthController) Login(c *gin.Context) {
	con.loginHandler(c)
}

// TODO add docs
func (con AuthController) RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register", nil)
}

// TODO add docs
func (con AuthController) RegisterView(c *gin.Context) {
	var newUser dto.RegisterDetails

	if err := c.BindJSON(&newUser); err != nil {
		// TODO
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := con.userService.Register(&newUser)
	if err != nil {
		// TODO
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "register-result", nil)
}

// TODO add docs
func (con AuthController) LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login", nil)
}

// TODO add docs
func (con AuthController) LoginView(c *gin.Context) {
	con.loginHandler(c)
	if c.Writer.Status() != http.StatusOK {
		// TODO
		return
	}
	// var loginVals dto.LoginDetails
	// c.BindJSON(&loginVals)

	// c.HTML(http.StatusOK, "login-result", gin.H{
	// 	"username": loginVals.Username,
	// })
}
