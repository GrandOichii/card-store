package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"store.api/auth"
	"store.api/config"
	"store.api/dto"
	"store.api/model"
	"store.api/query"
	"store.api/service"

	urlquery "github.com/google/go-querystring/query"
)

type CardController struct {
	config        *config.Configuration
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
		con.group.PATCH("/price/:id", con.UpdatePrice)
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

func NewCardController(config *config.Configuration, cardService service.CardService, auth gin.HandlerFunc, claimExtractF func(string, *gin.Context) (string, error)) *CardController {
	result := &CardController{
		config:        config,
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
// @Failure				400 {object} string
// @Failure				401 {object} string
// @Failure				403 {object} string
// @Router				/card [post]
func (con *CardController) Create(c *gin.Context) {
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

	var newCard dto.PostCard
	if err := c.BindJSON(&newCard); err != nil {
		AbortWithError(c, http.StatusBadRequest, err, true)
		return
	}

	card, err := con.cardService.Add(&newCard, uint(userId))
	if err != nil {
		AbortWithError(c, http.StatusBadRequest, err, true)
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
// @Failure				400 {object} string
// @Failure				404 {object} string
// @Router				/card/{id} [get]
func (con *CardController) ById(c *gin.Context) {
	p := c.Param("id")
	id, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		AbortWithError(c, http.StatusBadRequest, fmt.Errorf("%s is not a valid card id", p), true)
		return
	}

	card, err := con.cardService.GetById(uint(id))
	if err != nil {
		if err == service.ErrCardNotFound {
			AbortWithError(c, http.StatusNotFound, fmt.Errorf("no card with id %v", id), true)
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
// @Success				200 {object} service.CardQueryResult
// @Failure				400 {object} string
// @Router				/card [get]
func (con *CardController) Query(c *gin.Context) {
	var query query.CardQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		AbortWithError(c, http.StatusBadRequest, errors.New("invalid card query"), true)
		return
	}
	query.Keywords = strings.Join(strings.Fields(query.Keywords), " ")

	if len(strings.Split(query.Keywords, " ")) > int(con.config.Store.QueryKeywordLimit) {
		AbortWithError(c, http.StatusBadRequest, fmt.Errorf("too many keywords (limit: %d)", con.config.Store.QueryKeywordLimit), true)
		return
	}
	vals, err := urlquery.Values(query)
	if err != nil {
		// * should never happen
		panic(err)
	}
	query.Raw = vals.Encode()

	result := con.cardService.Query(&query)

	c.IndentedJSON(http.StatusOK, result)
}

// UpdateCard			godoc
// @Summary				Update card
// @Description			Updates an existing card
// @Param				Authorization header string false "Authenticator"
// @Param				id path int true "Card ID"
// @Param				card body dto.PostCard true "new card data"
// @Tags				Card
// @Success				200 {object} dto.GetCard
// @Failure				400 {object} string
// @Failure				401 {object} string
// @Failure				403 {object} string
// @Failure				404 {object} string
// @Router				/card/{id} [patch]
func (con *CardController) Update(c *gin.Context) {
	p := c.Param("id")
	id, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		AbortWithError(c, http.StatusBadRequest, fmt.Errorf("%s is not a valid card id", p), true)
		return
	}

	var newData dto.PostCard
	if err := c.BindJSON(&newData); err != nil {
		AbortWithError(c, http.StatusBadRequest, err, true)
		return
	}

	card, err := con.cardService.Update(&newData, uint(id))
	if err != nil {
		if err == service.ErrCardNotFound {
			AbortWithError(c, http.StatusNotFound, fmt.Errorf("no card with id %v", id), true)
			return
		}
		AbortWithError(c, http.StatusBadRequest, err, true)
		return
	}

	c.IndentedJSON(http.StatusOK, card)
}

// UpdatePrice			godoc
// @Summary				Update card price
// @Description			Updates an existing card's price
// @Param				Authorization header string false "Authenticator"
// @Param				id path int true "Card ID"
// @Param				price body dto.PriceUpdate true "new card price"
// @Tags				Card
// @Success				200 {object} dto.GetCard
// @Failure				400 {object} string
// @Failure				401 {object} string
// @Failure				403 {object} string
// @Failure				404 {object} string
// @Router				/card/price/{id} [patch]
func (con *CardController) UpdatePrice(c *gin.Context) {
	p := c.Param("id")
	id, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		AbortWithError(c, http.StatusBadRequest, fmt.Errorf("%s is not a valid card id", p), true)
		return
	}

	var newPrice dto.PriceUpdate
	if err := c.BindJSON(&newPrice); err != nil {
		AbortWithError(c, http.StatusBadRequest, err, true)
		return
	}

	card, err := con.cardService.UpdatePrice(uint(id), &newPrice)
	if err != nil {
		if err == service.ErrCardNotFound {
			AbortWithError(c, http.StatusNotFound, fmt.Errorf("no card with id %v", id), true)
			return
		}
		AbortWithError(c, http.StatusBadRequest, err, true)
		return
	}

	c.IndentedJSON(http.StatusOK, card)
}
