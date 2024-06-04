package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"store.api/auth"
	"store.api/dto"
	"store.api/model"
	"store.api/service"
)

type UserController struct {
	userService service.UserService
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
		con.group.GET("", con.GetInfo)
		con.group.GET("/login-test", func(ctx *gin.Context) {
			ctx.IndentedJSON(http.StatusOK, gin.H{
				"message": "hello:)",
			})
		})

		cart := con.group.Group("/cart")
		{
			cart.GET("", con.GetCart)
			cart.POST("", con.EditCartSlot)
		}
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

func NewUserController(userService service.UserService, cartService service.CartService, auth gin.HandlerFunc, claimExtractF func(string, *gin.Context) (string, error)) *UserController {
	return &UserController{
		userService:   userService,
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
// @Success				200 {object} dto.GetCart
// @Failure				401 {object} string
// @Router				/user/cart [get]
func (con *UserController) GetCart(c *gin.Context) {
	rawId, err := con.claimExtractF(auth.IDKey, c)
	if err != nil {
		AbortWithError(c, http.StatusUnauthorized, err, true)
		return
	}
	userId, err := strconv.ParseUint(rawId, 10, 32)
	if err != nil {
		AbortWithError(c, http.StatusUnauthorized, fmt.Errorf("%s is an invalid user id", rawId), true)
		return
	}

	cart, err := con.cartService.Get(uint(userId))
	if err != nil {
		if err == service.ErrUserNotFound {
			AbortWithError(c, http.StatusUnauthorized, fmt.Errorf("no user with id %d", userId), true)
			return
		}
		panic(err)
	}

	c.IndentedJSON(http.StatusOK, cart)
}

// EditCartSlot			godoc
// @Summary				Add, remove or alter cart slot
// @Description			Adds, removes or alters a cart slot
// @Param				Authorization header string false "Authenticator"
// @Param				collectionSlot body dto.PostCollectionSlot true "new cart slot data"
// @Tags				Collection
// @Success				200 {object} dto.GetCollection
// @Failure				400 {object} string
// @Failure				401 {object} string
// @Failure				404 {object} string
// @Router				/user/cart [post]
func (con *UserController) EditCartSlot(c *gin.Context) {
	rawId, err := con.claimExtractF(auth.IDKey, c)
	if err != nil {
		AbortWithError(c, http.StatusUnauthorized, err, true)
		return
	}
	userId, err := strconv.ParseUint(rawId, 10, 32)
	if err != nil {
		AbortWithError(c, http.StatusUnauthorized, fmt.Errorf("%s is an invalid user id", rawId), true)
		return
	}

	var newCartSlot dto.PostCartSlot
	if err := c.BindJSON(&newCartSlot); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := con.cartService.EditSlot(uint(userId), &newCartSlot)
	if err != nil {
		if err == service.ErrCardNotFound {
			AbortWithError(c, http.StatusNotFound, fmt.Errorf("no card with id %v", newCartSlot.CardId), true)
			return
		}
		AbortWithError(c, http.StatusBadRequest, err, true)
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}

// TODO add tests
// GetInfo				godoc
// @Summary				Get user info
// @Description			Gets the user's private information
// @Param				Authorization header string false "Authenticator"
// @Tags				User
// @Success				200 {object} dto.PrivateUserInfo
// @Failure				403 {object} string
// @Router				/user [get]
func (con *UserController) GetInfo(c *gin.Context) {
	rawId, err := con.claimExtractF(auth.IDKey, c)
	if err != nil {
		AbortWithError(c, http.StatusUnauthorized, err, true)
		return
	}

	userId, err := strconv.ParseUint(rawId, 10, 32)
	if err != nil {
		AbortWithError(c, http.StatusUnauthorized, fmt.Errorf("%s is an invalid user id", rawId), true)
		return
	}

	data, err := con.userService.ById(uint(userId))
	if err != nil {
		if err == service.ErrUserNotFound {
			// TODO repeated code
			AbortWithError(c, http.StatusUnauthorized, fmt.Errorf("no user with id %d", userId), true)
			return
		}
		panic(err)
	}

	c.IndentedJSON(http.StatusOK, data)
}
