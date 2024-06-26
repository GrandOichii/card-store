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

func newUserController(userService service.UserService, cartService service.CartService) *controller.UserController {
	return controller.NewUserController(
		userService,
		cartService,
		func(ctx *gin.Context) {},
		func(s string, ctx *gin.Context) (string, error) {
			return "1", nil
		},
	)
}

func Test_User_ShouldGetCart(t *testing.T) {
	// arrange
	userService := newMockUserService()
	cartService := newMockCartService()
	controller := newUserController(userService, cartService)
	cartService.On("Get", mock.Anything).Return(&dto.GetCart{}, nil)
	c, w := createTestContext(nil)

	// act
	controller.GetCart(c)

	// assert
	assert.Equal(t, 200, w.Code)
}

func Test_User_ShouldNotGetCartNoUser(t *testing.T) {
	// arrange
	userService := newMockUserService()
	cartService := newMockCartService()
	controller := newUserController(userService, cartService)
	cartService.On("Get", mock.Anything).Return(nil, service.ErrUserNotFound)
	c, w := createTestContext(nil)

	// act
	controller.GetCart(c)

	// assert
	assert.Equal(t, 401, w.Code)
}

func Test_User_ShouldEditCartSlot(t *testing.T) {
	// arrange
	userService := newMockUserService()
	cartService := newMockCartService()
	controller := newUserController(userService, cartService)
	cartService.On("EditSlot", mock.Anything, mock.Anything).Return(&dto.GetCart{}, nil)
	c, w := createTestContext(&dto.PostCartSlot{
		CardId: 1,
		Amount: 1,
	})

	// act
	controller.EditCartSlot(c)

	// assert
	assert.Equal(t, 200, w.Code)
}

func Test_User_ShouldNotEditCartSlotCardNotFound(t *testing.T) {
	// arrange
	userService := newMockUserService()
	cartService := newMockCartService()
	controller := newUserController(userService, cartService)
	cartService.On("EditSlot", mock.Anything, mock.Anything).Return(nil, service.ErrCardNotFound)
	c, w := createTestContext(&dto.PostCartSlot{
		CardId: 1,
		Amount: 1,
	})

	// act
	controller.EditCartSlot(c)

	// assert
	assert.Equal(t, 404, w.Code)
}

func Test_User_ShouldNotEditCartSlotBadRequest(t *testing.T) {
	// arrange
	userService := newMockUserService()
	cartService := newMockCartService()
	controller := newUserController(userService, cartService)
	cartService.On("EditSlot", mock.Anything, mock.Anything).Return(nil, errors.New(""))
	c, w := createTestContext(&dto.PostCartSlot{
		CardId: 1,
		Amount: 1,
	})

	// act
	controller.EditCartSlot(c)

	// assert
	assert.Equal(t, 400, w.Code)
}

func Test_User_ShouldGetInfo(t *testing.T) {
	// arrange
	userService := newMockUserService()
	cartService := newMockCartService()
	controller := newUserController(userService, cartService)
	userService.On("ById", mock.Anything).Return(&dto.PrivateUserInfo{}, nil)
	c, w := createTestContext(&dto.PostCartSlot{
		CardId: 1,
		Amount: 1,
	})

	// act
	controller.GetInfo(c)

	// assert
	assert.Equal(t, 200, w.Code)
}

func Test_User_ShouldNotGetInfo(t *testing.T) {
	// arrange
	userService := newMockUserService()
	cartService := newMockCartService()
	controller := newUserController(userService, cartService)
	userService.On("ById", mock.Anything).Return(nil, service.ErrUserNotFound)
	c, w := createTestContext(&dto.PostCartSlot{
		CardId: 1,
		Amount: 1,
	})

	// act
	controller.GetInfo(c)

	// assert
	assert.Equal(t, 401, w.Code)
}
