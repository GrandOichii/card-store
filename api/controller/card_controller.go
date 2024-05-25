package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"store.api/auth"
	"store.api/dto"
	"store.api/model"
	"store.api/query"
	"store.api/service"
)

type CardController struct {
	cardService   service.CardService
	auth          gin.HandlerFunc
	claimExtractF func(string, *gin.Context) (string, error)

	group       *gin.RouterGroup
	authChecker auth.AuthorizationChecker
}

func (con *CardController) ConfigureApi(r *gin.RouterGroup) {
	r.GET("/card", con.Query)
	r.GET("/card/:id", con.ById)
	con.group = r.Group("/card")
	{
		con.group.Use(con.auth)
		con.group.POST("", con.Create)
		con.group.PATCH("/:id", con.Update)
	}

	path := con.group.BasePath() + "*"
	con.authChecker = auth.NewAuthorizationCheckerBuilder().
		ForPath(path).
		ForAnyMethod().
		PermitAll().
		ForPath(path).
		ForMethod("POST").
		Permit(func(user *model.User) bool {
			return user.IsAdmin && user.Verified
		}).
		ForPath(path).
		ForMethod("PATCH").
		Permit(func(user *model.User) bool {
			return user.IsAdmin && user.Verified
		}).
		Build()
}

func (con *CardController) Check(c *gin.Context, user *model.User) (authorized bool, matches bool) {
	return con.authChecker.Check(c, user)
}

func NewCardController(cardService service.CardService, auth gin.HandlerFunc, claimExtractF func(string, *gin.Context) (string, error)) *CardController {
	result := &CardController{
		cardService:   cardService,
		auth:          auth,
		claimExtractF: claimExtractF,
	}

	return result
}

// CreateCard			godoc
// @Summary				Create new card
// @Description			Creates a new card
// @Param				Authorization header string false "Authenticator"
// @Param				card body dto.PostCard true "new card data"
// @Tags				Card
// @Success				201 {object} dto.GetCard
// @Failure				400 {object} ErrResponse
// @Failure				401 {object} ErrResponse
// @Failure				403 {object} ErrResponse
// @Router				/card [post]
func (con *CardController) Create(c *gin.Context) {
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

	var newCard dto.PostCard
	if err := c.BindJSON(&newCard); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	card, err := con.cardService.Add(&newCard, uint(userId))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, card)
}

// ById					godoc
// @Summary				Fetch card by id
// @Description			Fetches a card by it's id
// @Param				id path int true "Card ID"
// @Tags				Card
// @Success				200 {object} dto.GetCard
// @Failure				400 {object} ErrResponse
// @Failure				404 {object} ErrResponse
// @Router				/card/{id} [get]
func (con *CardController) ById(c *gin.Context) {
	p := c.Param("id")
	id, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("%s is not a valid card id", p),
		})
		return
	}

	card, err := con.cardService.GetById(uint(id))
	if err != nil {
		if err == service.ErrCardNotFound {
			c.IndentedJSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf("no card with id %v", id),
			})
			return
		}
		panic(err)
	}

	c.IndentedJSON(http.StatusOK, card)
}

// Query				godoc
// @Summary				Fetch card by query
// @Description			Fetches all cards that match the query
// @Param				query query query.CardQuery false "Card query"
// @Tags				Card
// @Success				200 {object} dto.GetCard[]
// @Failure				400 {object} ErrResponse
// @Router				/card [get]
func (con *CardController) Query(c *gin.Context) {
	var query query.CardQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": "invalid card query",
		})
		return
	}

	result := con.cardService.Query(&query)

	c.IndentedJSON(http.StatusOK, result)
}

// CreateCard			godoc
// @Summary				Update card
// @Description			Updates an existing card
// @Param				Authorization header string false "Authenticator"
// @Param				id path int true "Card ID"
// @Param				card body dto.PostCard true "new card data"
// @Tags				Card
// @Success				200 {object} dto.GetCard
// @Failure				400 {object} ErrResponse
// @Failure				401 {object} ErrResponse
// @Failure				403 {object} ErrResponse
// @Failure				404 {object} ErrResponse
// @Router				/card/:id [patch]
func (con *CardController) Update(c *gin.Context) {
	p := c.Param("id")
	id, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("%s is not a valid card id", p),
		})
		return
	}

	var newData dto.PostCard
	if err := c.BindJSON(&newData); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	card, err := con.cardService.Update(&newData, uint(id))
	if err != nil {
		if err == service.ErrCardNotFound {
			c.IndentedJSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf("no card with id %v", id),
			})
			return
		}
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, card)
}
