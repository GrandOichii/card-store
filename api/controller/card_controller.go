package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"store.api/auth"
	"store.api/dto"
	"store.api/service"
)

type CardController struct {
	Controller

	cardService   service.CardService
	loginHandler  gin.HandlerFunc
	claimExtractF func(string, *gin.Context) (string, error)
}

func (con CardController) Configure(r *gin.RouterGroup) {
	// TODO remove
	r.GET("/card", con.All)

	g := r.Group("/card")
	{
		g.Use(con.loginHandler)
		g.POST("", con.Create)
	}
}

func NewCardController(cardService service.CardService, loginHandler gin.HandlerFunc, claimExtractF func(string, *gin.Context) (string, error)) *CardController {
	return &CardController{
		cardService:   cardService,
		loginHandler:  loginHandler,
		claimExtractF: claimExtractF,
	}
}

// TODO add other @Failure docs
// example: @Failure      400  {object}  httputil.HTTPError

// AllCards				godoc
// @Summary				Fetch all cards
// @Description			Fetches all existing cards
// @Tags				Card
// @Success				200 {object} dto.GetCard[]
// @Router				/card [get]
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

// // UserRegister			godoc
// // @Summary				Registers the user
// // @Description			Checks the user data and adds it to the repo
// // @Param				details body dto.RegisterDetails true "Register details"
// // @Tags				Auth
// // @Success				201
// // @Router				/auth/register [post]
// func (con *CardController) Register(c *gin.Context) {
// 	var newUser dto.RegisterDetails

// 	if err := c.BindJSON(&newUser); err != nil {
// 		c.IndentedJSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	err := con.cardService.Register(&newUser)
// 	if err != nil {
// 		c.IndentedJSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}
// 	c.Status(http.StatusCreated)
// }
