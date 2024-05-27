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

func newCardController(cardService service.CardService) *controller.CardController {
	return controller.NewCardController(
		cardService,
		func(*gin.Context) {},
		func(s string, ctx *gin.Context) (string, error) {
			return "1", nil
		},
		// auth.NewJwtMiddleware(&config.Configuration{
		// 	AuthKey: "test secret key",
		// }, authService, repo).Middle.LoginHandler,
	)
}

func Test_Card_ShouldCreate(t *testing.T) {
	// arrange
	service := newMockCardService()
	controller := newCardController(service)
	service.On("Add", mock.Anything, mock.Anything).Return(&dto.GetCard{}, nil)
	data := dto.PostCard{
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

func Test_Card_ShouldNotCreate(t *testing.T) {
	// arrange
	service := newMockCardService()
	controller := newCardController(service)
	service.On("Add", mock.Anything, mock.Anything).Return(nil, errors.New(""))
	data := dto.PostCard{
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

func Test_Card_ShouldNotCreateBadData(t *testing.T) {
	// arrange
	service := newMockCardService()
	controller := newCardController(service)
	service.On("Add", mock.Anything, mock.Anything).Return(nil, errors.New(""))
	data := []string{"first", "second"}
	c, w := createTestContext(data)

	// act
	controller.Create(c)

	// assert
	assert.Equal(t, 400, w.Code)
}

func Test_Card_ShouldFetchById(t *testing.T) {
	// arrange
	service := newMockCardService()
	controller := newCardController(service)
	service.On("GetById", mock.Anything).Return(&dto.GetCard{}, nil)
	c, w := createTestContext(nil)
	c.AddParam("id", "12")

	// act
	controller.ById(c)

	// assert
	assert.Equal(t, 200, w.Code)
}

func Test_Card_ShouldNotFetchById(t *testing.T) {
	// arrange
	s := newMockCardService()
	controller := newCardController(s)
	s.On("GetById", mock.Anything).Return(nil, service.ErrCardNotFound)
	c, w := createTestContext(nil)
	c.AddParam("id", "1")

	// act
	controller.ById(c)

	// assert
	assert.Equal(t, 404, w.Code)
}

func Test_Card_ShouldFetchByType(t *testing.T) {
	// arrange
	service := newMockCardService()
	controller := newCardController(service)
	service.On("Query", mock.Anything).Return([]*dto.GetCard{})
	c, w := createTestContext(nil)
	c.Request.URL, _ = url.Parse("?type=CT1")

	// act
	controller.Query(c)

	// assert
	assert.Equal(t, 200, w.Code)
}

func Test_Card_ShouldFetchByName(t *testing.T) {
	// arrange
	service := newMockCardService()
	controller := newCardController(service)
	service.On("Query", mock.Anything).Return([]*dto.GetCard{})
	c, w := createTestContext(nil)
	c.Request.URL, _ = url.Parse("?name=card")

	// act
	controller.Query(c)

	// assert
	assert.Equal(t, 200, w.Code)
}

func Test_Card_ShouldFetchByMinPrice(t *testing.T) {
	// arrange
	service := newMockCardService()
	controller := newCardController(service)
	service.On("Query", mock.Anything).Return([]*dto.GetCard{})
	c, w := createTestContext(nil)
	c.Request.URL, _ = url.Parse("?minPrice=30")

	// act
	controller.Query(c)

	// assert
	assert.Equal(t, 200, w.Code)
}

func Test_Card_ShouldFetchByMaxPrice(t *testing.T) {
	// arrange
	service := newMockCardService()
	controller := newCardController(service)
	service.On("Query", mock.Anything).Return([]*dto.GetCard{})
	c, w := createTestContext(nil)
	c.Request.URL, _ = url.Parse("?maxPrice=400")

	// act
	controller.Query(c)

	// assert
	assert.Equal(t, 200, w.Code)
}

func Test_Card_ShouldUpdate(t *testing.T) {
	// arrange
	service := newMockCardService()
	controller := newCardController(service)
	service.On("Update", mock.Anything, mock.Anything).Return(&dto.GetCard{}, nil)
	data := dto.PostCard{
		Name:     "card name",
		Text:     "card text",
		Price:    10,
		Type:     "CT1",
		Language: "ENG",
	}
	c, w := createTestContext(data)
	c.AddParam("id", "12")

	// act
	controller.Update(c)

	// assert
	assert.Equal(t, 200, w.Code)
}

func Test_Card_ShouldNotUpdateCollectionNotFound(t *testing.T) {
	// arrange
	s := newMockCardService()
	controller := newCardController(s)
	s.On("Update", mock.Anything, mock.Anything).Return(nil, service.ErrCardNotFound)
	data := dto.PostCard{
		Name:     "card name",
		Text:     "card text",
		Price:    10,
		Type:     "CT1",
		Language: "ENG",
	}
	c, w := createTestContext(data)
	c.AddParam("id", "12")

	// act
	controller.Update(c)

	// assert
	assert.Equal(t, 404, w.Code)
}

func Test_Card_ShouldNotUpdateBadRequest(t *testing.T) {
	// arrange
	s := newMockCardService()
	controller := newCardController(s)
	s.On("Update", mock.Anything, mock.Anything).Return(nil, errors.New(""))
	data := dto.PostCard{
		Name:     "card name",
		Text:     "card text",
		Price:    10,
		Type:     "CT1",
		Language: "ENG",
	}
	c, w := createTestContext(data)
	c.AddParam("id", "12")

	// act
	controller.Update(c)

	// assert
	assert.Equal(t, 400, w.Code)
}
