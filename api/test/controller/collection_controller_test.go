package controller_test

import (
	"github.com/gin-gonic/gin"
	"store.api/controller"
	"store.api/service"
)

func createCollectionController(collectionService service.CollectionService) *controller.CollectionController {
	return controller.NewCollectionController(
		collectionService,
		func(*gin.Context) {},
		func(s string, ctx *gin.Context) (string, error) {
			return "1", nil
		},
		// auth.NewJwtMiddleware(&config.Configuration{
		// 	AuthKey: "test secret key",
		// }, userService, repo).Middle.LoginHandler,
	)
}
