package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"store.api/auth"
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
		// TODO
	}

	con.authChecker = auth.NewAuthorizationCheckerBuilder().
		ForPath(con.group.BasePath() + "/*").
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
