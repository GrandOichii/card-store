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
		con.group.GET("/:id", con.ById)
		con.group.POST("", con.Create)
		con.group.POST("/:collectionId", con.EditSlot)
		con.group.DELETE("/:id", con.Delete)
		con.group.PATCH("/:id", con.UpdateInfo)
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
// @Failure				401 {object} string
// @Router				/collection/all [get]
func (con *CollectionController) All(c *gin.Context) {
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

	collections := con.collectionService.GetAll(uint(userId))

	c.IndentedJSON(http.StatusOK, collections)
}

// PostCollection		godoc
// @Summary				Create new collection
// @Description			Creates a new card collection
// @Param				Authorization header string false "Authenticator"
// @Param				collection body dto.PostCollection true "new collection data"
// @Tags				Collection
// @Success				201 {object} dto.GetCollection
// @Failure				400 {object} string
// @Failure				401 {object} string
// @Router				/collection [post]
func (con *CollectionController) Create(c *gin.Context) {
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

	var newCollection dto.PostCollection
	if err := c.BindJSON(&newCollection); err != nil {
		AbortWithError(c, http.StatusBadRequest, err, true)
		return
	}

	collection, err := con.collectionService.Create(&newCollection, uint(userId))
	if err != nil {
		if err == service.ErrNotVerified {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		if err == service.ErrUserNotFound {
			AbortWithError(c, http.StatusUnauthorized, userNotFound(uint(userId)), true)
			return
		}
		AbortWithError(c, http.StatusBadRequest, err, true)
		return
	}

	c.IndentedJSON(http.StatusCreated, collection)
}

// EditSlot				godoc
// @Summary				Add, remove or alter collection slot
// @Description			Adds, removes or alters a collection slot in an existing collection
// @Param				Authorization header string false "Authenticator"
// @Param				collectionId path int true "Collection ID"
// @Param				collectionSlot body dto.PostCollectionSlot true "new collection slot data"
// @Tags				Collection
// @Success				201 {object} dto.GetCollection
// @Failure				400 {object} string
// @Failure				401 {object} string
// @Failure				404 {object} string
// @Router				/collection/{collectionId} [post]
func (con *CollectionController) EditSlot(c *gin.Context) {
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

	p := c.Param("collectionId")
	collectionId, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		AbortWithError(c, http.StatusBadRequest, fmt.Errorf("%s is not a valid collection id", p), true)
		return
	}

	var newColSlot dto.PostCollectionSlot
	if err := c.BindJSON(&newColSlot); err != nil {
		AbortWithError(c, http.StatusBadRequest, err, true)
		return
	}

	collection, err := con.collectionService.EditSlot(&newColSlot, uint(collectionId), uint(userId))
	if err != nil {
		if err == service.ErrNotVerified {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		if err == service.ErrCardNotFound {
			AbortWithError(c, http.StatusNotFound, fmt.Errorf("no card with id %d", newColSlot.CardId), true)
			return
		}
		if err == service.ErrUserNotFound {
			AbortWithError(c, http.StatusUnauthorized, userNotFound(uint(userId)), true)
			return
		}
		AbortWithError(c, http.StatusBadRequest, err, true)
		return
	}

	c.IndentedJSON(http.StatusOK, collection)
}

// ById					godoc
// @Summary				Fetch collection by id
// @Description			Fetches a collection by it's id
// @Param				id path int true "Collection ID"
// @Tags				Collection
// @Param				Authorization header string false "Authenticator"
// @Success				200 {object} dto.GetCollection
// @Failure				400 {object} string
// @Failure				401 {object} string
// @Failure				404 {object} string
// @Router				/collection/{id} [get]
func (con *CollectionController) ById(c *gin.Context) {
	p := c.Param("id")
	id, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		AbortWithError(c, http.StatusBadRequest, fmt.Errorf("%s is not a valid card id", p), true)
		return
	}

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

	collection, err := con.collectionService.GetById(uint(id), uint(userId))
	if err != nil {
		if err == service.ErrCollectionNotFound {
			AbortWithError(c, http.StatusNotFound, fmt.Errorf("no collection with id %d", id), true)
			return
		}
		panic(err)
	}

	c.IndentedJSON(http.StatusOK, collection)
}

// ById					godoc
// @Summary				Delete collection
// @Description			Deletes a collection by it's id
// @Param				id path int true "Collection ID"
// @Param				Authorization header string false "Authenticator"
// @Tags				Collection
// @Success				200
// @Failure				400 {object} string
// @Failure				401 {object} string
// @Failure				404 {object} string
// @Router				/collection/{id} [delete]
func (con *CollectionController) Delete(c *gin.Context) {
	p := c.Param("id")
	id, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		AbortWithError(c, http.StatusBadRequest, fmt.Errorf("%s is not a valid card id", p), true)
		return
	}

	rawId, err := con.claimExtractF(auth.IDKey, c)
	if err != nil {
		AbortWithError(c, http.StatusUnauthorized, err, true)
		return
	}
	userId, err := strconv.ParseUint(rawId, 10, 32)
	if err != nil {
		AbortWithError(c, http.StatusUnauthorized, fmt.Errorf("%s is an invalid collection id", rawId), true)
		return
	}

	err = con.collectionService.Delete(uint(id), uint(userId))
	if err != nil {
		AbortWithError(c, http.StatusNotFound, err, true)
		return
	}

	c.Status(http.StatusOK)
}

// UpdateInfo			godoc
// @Summary				Update collection info
// @Description			Deletes a collection's info by it's id
// @Param				id path int true "Collection ID"
// @Param				Authorization header string false "Authenticator"
// @Param				collection body dto.PostCollection true "new collection data"
// @Tags				Collection
// @Success				200
// @Failure				400 {object} string
// @Failure				401 {object} string
// @Failure				404 {object} string
// @Router				/collection/{id} [patch]
func (con *CollectionController) UpdateInfo(c *gin.Context) {
	p := c.Param("id")
	id, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		AbortWithError(c, http.StatusBadRequest, fmt.Errorf("%s is not a valid card id", p), true)
		return
	}

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

	var newData dto.PostCollection
	if err := c.BindJSON(&newData); err != nil {
		AbortWithError(c, http.StatusBadRequest, err, true)
		return
	}

	collection, err := con.collectionService.UpdateInfo(&newData, uint(id), uint(userId))
	if err != nil {
		if err == service.ErrNotVerified {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		if err == service.ErrCollectionNotFound {
			AbortWithError(c, http.StatusNotFound, fmt.Errorf("no collection with id %d", id), true)
			return
		}
		if err == service.ErrUserNotFound {
			AbortWithError(c, http.StatusNotFound, userNotFound(uint(userId)), true)
			return
		}
		AbortWithError(c, http.StatusBadRequest, err, true)
		return
	}

	c.IndentedJSON(http.StatusOK, collection)
}
