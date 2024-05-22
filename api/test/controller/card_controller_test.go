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

func createCardController(cardService service.CardService) *controller.CardController {
	return controller.NewCardController(
		cardService,
		func(*gin.Context) {},
		func(s string, ctx *gin.Context) (string, error) {
			return "userID", nil
		},
		// auth.NewJwtMiddleware(&config.Configuration{
		// 	AuthKey: "test secret key",
		// }, userService, repo).Middle.LoginHandler,
	)
}

func Test_ShouldFetchAll(t *testing.T) {
	// arrange
	service := createMockCardService()
	controller := createCardController(service)
	service.On("GetAll").Return([]*dto.GetCard{})

	c, w := createTestContext(nil)

	// act
	controller.All(c)

	// assert
	assert.Equal(t, w.Code, 200)
}

func Test_ShouldCreate(t *testing.T) {
	// arrange
	service := createMockCardService()
	controller := createCardController(service)
	service.On("Add", mock.Anything, mock.Anything).Return(&dto.GetCard{}, nil)
	data := dto.CreateCard{
		Name:  "card name",
		Text:  "card text",
		Price: 10,
	}
	c, w := createTestContext(data)

	// act
	controller.Create(c)

	// assert
	assert.Equal(t, w.Code, 201)
}

func Test_ShouldNotCreate(t *testing.T) {
	// arrange
	service := createMockCardService()
	controller := createCardController(service)
	service.On("Add", mock.Anything, mock.Anything).Return(nil, errors.New(""))
	data := dto.CreateCard{
		Name:  "card name",
		Text:  "card text",
		Price: 10,
	}
	c, w := createTestContext(data)

	// act
	controller.Create(c)

	// assert
	assert.Equal(t, w.Code, 400)
}

func Test_ShouldNotCreateBadData(t *testing.T) {
	// arrange
	service := createMockCardService()
	controller := createCardController(service)
	service.On("Add", mock.Anything, mock.Anything).Return(nil, errors.New(""))
	data := []string{"first", "second"}
	c, w := createTestContext(data)

	// act
	controller.Create(c)

	// assert
	assert.Equal(t, w.Code, 400)
}

func Test_ShouldFetchById(t *testing.T) {
	// arrange
	service := createMockCardService()
	controller := createCardController(service)
	service.On("GetById", mock.Anything).Return(&dto.GetCard{}, nil)
	c, w := createTestContext(nil)
	c.AddParam("id", "12")

	// act
	controller.ById(c)

	// assert
	assert.Equal(t, w.Code, 200)
}

func Test_ShouldNotFetchById(t *testing.T) {
	// arrange
	service := createMockCardService()
	controller := createCardController(service)
	service.On("GetById", mock.Anything).Return(nil, errors.New(""))
	c, w := createTestContext(nil)
	c.AddParam("id", "1")

	// act
	controller.ById(c)

	// assert
	assert.Equal(t, w.Code, 404)
}
