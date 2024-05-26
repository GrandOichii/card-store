package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"store.api/auth"
	"store.api/model"
	"store.api/service"
)

type UserController struct {
	cartService service.CartService

	group         *gin.RouterGroup
	auth          gin.HandlerFunc
	authChecker   auth.AuthorizationChecker
	claimExtractF func(string, *gin.Context) (string, error)
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

		cart := con.group.Group("/cart")
		{
			cart.GET("", con.GetCart)
		}
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

func NewUserController(cartService service.CartService, auth gin.HandlerFunc, claimExtractF func(string, *gin.Context) (string, error)) *UserController {
	return &UserController{
		cartService:   cartService,
		auth:          auth,
		claimExtractF: claimExtractF,
	}
}

// GetCart				godoc
// @Summary				Fetch cart
// @Description			Fetches the user's cart
// @Param				Authorization header string false "Authenticator"
// @Tags				Cart
// @Success				200 {object} dto.GetCard
// @Failure				401 {object} ErrResponse
// @Router				/user/cart [get]
func (con *UserController) GetCart(c *gin.Context) {
	rawId, err := con.claimExtractF(auth.IDKey, c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	userId, err := strconv.ParseUint(rawId, 10, 32)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("%s is an invalid user id", rawId))
		return
	}

	cart, err := con.cartService.Get(uint(userId))
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	c.IndentedJSON(http.StatusOK, cart)
}
