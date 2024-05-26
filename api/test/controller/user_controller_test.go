package controller_test

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"store.api/controller"
	"store.api/dto"
	"store.api/service"
)

func newUserController(cartService service.CartService) *controller.UserController {
	return controller.NewUserController(
		cartService,
		func(ctx *gin.Context) {},
		func(s string, ctx *gin.Context) (string, error) {
			return "1", nil
		},
	)
}

func Test_User_ShouldGetCart(t *testing.T) {
	// arrange
	cartService := newMockCartService()
	controller := newUserController(cartService)
	cartService.On("Get", mock.Anything).Return(&dto.GetCart{}, nil)
	c, w := createTestContext(nil)

	// act
	controller.GetCart(c)

	// assert
	assert.Equal(t, 200, w.Code)
}

func Test_User_ShouldNotGetCartNoUser(t *testing.T) {
	// arrange
	cartService := newMockCartService()
	controller := newUserController(cartService)
	cartService.On("Get", mock.Anything).Return(nil, errors.New(""))
	c, w := createTestContext(nil)

	// act
	controller.GetCart(c)

	// assert
	assert.Equal(t, 401, w.Code)
}
