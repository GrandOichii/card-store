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

type CardController struct {
	cardService   service.CardService
	loginHandler  gin.HandlerFunc
	claimExtractF func(string, *gin.Context) (string, error)

	group       *gin.RouterGroup
	authChecker auth.AuthorizationChecker
}

func (con *CardController) ConfigureApi(r *gin.RouterGroup) {
	// TODO remove
	r.GET("/card/all", con.All)

	r.GET("/card/:id", con.ById)
	con.group = r.Group("/card")
	{
		con.group.Use(con.loginHandler)
		con.group.POST("", con.Create)
	}

	con.authChecker = auth.NewAuthorizationCheckerBuilder().
		ForPath(con.group.BasePath()).
		ForAnyMethod().
		PermitAll().
		ForPath(con.group.BasePath()).
		ForMethod("POST").
		Permit(func(user *model.User) bool {
			return user.IsAdmin && user.Verified
		}).
		Build()
}

func (con *CardController) ConfigurePages(r *gin.RouterGroup) {
	r.GET("view/card/id-search", func(c *gin.Context) {
		c.HTML(http.StatusOK, "card-id-search.html", nil)
	})
	r.GET("view/card/all", func(c *gin.Context) {
		cards := con.cardService.GetAll()
		c.HTML(http.StatusOK, "card-list", gin.H{
			"cards": cards,
		})
	})
	r.GET("view/card", func(c *gin.Context) {
		p := c.Query("id")
		id, err := strconv.ParseUint(p, 10, 32)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		card, err := con.cardService.GetById(uint(id))
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		c.HTML(http.StatusOK, "card", card)
	})

}

func (con *CardController) Check(c *gin.Context, user *model.User) (authorized bool, matches bool) {
	return con.authChecker.Check(c, user)
}

func NewCardController(cardService service.CardService, loginHandler gin.HandlerFunc, claimExtractF func(string, *gin.Context) (string, error)) *CardController {
	result := &CardController{
		cardService:   cardService,
		loginHandler:  loginHandler,
		claimExtractF: claimExtractF,
	}

	return result
}

// TODO add other @Failure docs
// example: @Failure      400  {object}  httputil.HTTPError

// AllCards				godoc
// @Summary				Fetch all cards
// @Description			Fetches all existing cards
// @Tags				Card
// @Success				200 {object} dto.GetCard[]
// @Router				/card/all [get]
func (con *CardController) All(c *gin.Context) {
	cards := con.cardService.GetAll()
	c.IndentedJSON(http.StatusOK, cards)
}

// CreateCard			godoc
// @Summary				Create new card
// @Description			Creates a new card
// @Param				Authorization header string true "Authenticator"
// @Param				card body dto.CreateCard true "new card data"
// @Tags				Card
// @Success				201 {object} dto.GetCard
// @Router				/card [post]
func (con *CardController) Create(c *gin.Context) {
	// TODO! make admin only

	username, err := con.claimExtractF(auth.IDKey, c)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	var newCard dto.CreateCard
	if err := c.BindJSON(&newCard); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	card, err := con.cardService.Add(&newCard, username)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, card)
}

// AllCards				godoc
// @Summary				Fetch card by id
// @Description			Fetches a card by it's id
// @Param				id path int true "Card ID"
// @Tags				Card
// @Success				200 {object} dto.GetCard
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
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, card)
}
