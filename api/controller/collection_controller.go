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

type CollectionController struct {
	collectionService service.CollectionService
	auth              gin.HandlerFunc
	claimExtractF     func(string, *gin.Context) (string, error)

	group       *gin.RouterGroup
	authChecker auth.AuthorizationChecker
}

func (con *CollectionController) ConfigureApi(r *gin.RouterGroup) {
	con.group = r.Group("/collection")
	con.group.Use(con.auth)
	{
		con.group.GET("/all", con.All)
		con.group.POST("", con.Create)
		con.group.POST("/:collectionId", con.AddCard)
		// POST: add card
		// PUT: modify card amount (can delete card)
		// DELETE: remove collection
	}

	con.authChecker = auth.NewAuthorizationCheckerBuilder().
		ForPath(con.group.BasePath() + "*").
		ForMethod("*").
		PermitAll().
		Build()
}

func (con *CollectionController) Check(c *gin.Context, user *model.User) (authorized bool, matches bool) {
	return con.authChecker.Check(c, user)
}

func NewCollectionController(collectionService service.CollectionService, auth gin.HandlerFunc, claimExtractF func(string, *gin.Context) (string, error)) *CollectionController {
	return &CollectionController{
		collectionService: collectionService,
		auth:              auth,
		claimExtractF:     claimExtractF,
	}
}

// All					godoc
// @Summary				Fetch all collections
// @Description			Fetches all the user's collections
// @Param				Authorization header string false "Authenticator"
// @Tags				Collection
// @Success				200 {object} dto.GetCollection[]
// @Failure				401 {object} ErrResponse
// @Router				/collection/all [get]
func (con *CollectionController) All(c *gin.Context) {
	rawId, err := con.claimExtractF(auth.IDKey, c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}
	userId, err := strconv.ParseUint(rawId, 10, 32)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	collections := con.collectionService.GetAll(uint(userId))

	c.IndentedJSON(http.StatusOK, collections)
}

// CreateCollection		godoc
// @Summary				Create new collection
// @Description			Creates a new card collection
// @Param				Authorization header string false "Authenticator"
// @Param				collection body dto.CreateCollection true "new collection data"
// @Tags				Collection
// @Success				201 {object} dto.GetCollection
// @Failure				400 {object} ErrResponse
// @Failure				401
// @Router				/collection [post]
func (con *CollectionController) Create(c *gin.Context) {
	rawId, err := con.claimExtractF(auth.IDKey, c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}
	userId, err := strconv.ParseUint(rawId, 10, 32)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	var newCollection dto.CreateCollection
	if err := c.BindJSON(&newCollection); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	collection, err := con.collectionService.Create(&newCollection, uint(userId))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, collection)
}

// AddCard				godoc
// @Summary				Add new card slot
// @Description			Adds a new card slot to an existing collection
// @Param				Authorization header string false "Authenticator"
// @Param				collectionId path int true "Collection ID"
// @Param				cardSlot body dto.CreateCardSlot true "new card slot data"
// @Tags				Collection
// @Success				201 {object} dto.GetCollection
// @Failure				400 {object} ErrResponse
// @Failure				401
// @Router				/collection/{collectionId} [post]
func (con *CollectionController) AddCard(c *gin.Context) {
	rawId, err := con.claimExtractF(auth.IDKey, c)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}
	userId, err := strconv.ParseUint(rawId, 10, 32)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	p := c.Param("collectionId")
	collectionId, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("%s is not a valid collection id", p),
		})
		return
	}

	var newCardSlot dto.CreateCardSlot
	if err := c.BindJSON(&newCardSlot); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	collection, err := con.collectionService.AddCard(&newCardSlot, uint(collectionId), uint(userId))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, collection)
}
