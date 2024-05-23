package controller_test

import (
	"errors"
	"net/url"
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

func Test_ShouldCreate(t *testing.T) {
	// arrange
	service := createMockCardService()
	controller := createCardController(service)
	service.On("Add", mock.Anything, mock.Anything).Return(&dto.GetCard{}, nil)
	data := dto.CreateCard{
		Name:  "card name",
		Text:  "card text",
		Price: 10,
		Type:  "CT1",
	}
	c, w := createTestContext(data)

	// act
	controller.Create(c)

	// assert
	assert.Equal(t, 201, w.Code)
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
		Type:  "CT1",
	}
	c, w := createTestContext(data)

	// act
	controller.Create(c)

	// assert
	assert.Equal(t, 400, w.Code)
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
	assert.Equal(t, 400, w.Code)
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
	assert.Equal(t, 200, w.Code)
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
	assert.Equal(t, 404, w.Code)
}

// TODO replace when querying will be more complex
func Test_ShouldFetchByType(t *testing.T) {
	// arrange
	service := createMockCardService()
	controller := createCardController(service)
	service.On("GetByType", mock.Anything).Return([]*dto.GetCard{})
	c, w := createTestContext(nil)
	c.Request.URL, _ = url.Parse("?type=CT1")

	// act
	controller.Query(c)

	// assert
	assert.Equal(t, 200, w.Code)
}
