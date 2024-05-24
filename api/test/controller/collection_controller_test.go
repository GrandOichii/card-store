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

func Test_Collection_ShouldCreate(t *testing.T) {
	// arrange
	service := createMockCollectionService()
	controller := createCollectionController(service)
	service.On("Create", mock.Anything, mock.Anything).Return(&dto.GetCollection{}, nil)
	data := dto.CreateCollection{
		Name:        "collection1",
		Description: "collection description",
	}
	c, w := createTestContext(data)

	// act
	controller.Create(c)

	// assert
	assert.Equal(t, 201, w.Code)
}

func Test_Collection_ShouldNotCreate(t *testing.T) {
	// arrange
	service := createMockCollectionService()
	controller := createCollectionController(service)
	service.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New(""))
	data := dto.CreateCollection{
		Name:        "collection1",
		Description: "collection description",
	}
	c, w := createTestContext(data)

	// act
	controller.Create(c)

	// assert
	assert.Equal(t, 400, w.Code)
}

func Test_Collection_ShouldNotCreateUnverified(t *testing.T) {
	// arrange
	s := createMockCollectionService()
	controller := createCollectionController(s)
	s.On("Create", mock.Anything, mock.Anything).Return(nil, service.ErrNotVerified)
	data := dto.CreateCollection{
		Name:        "collection1",
		Description: "collection description",
	}
	c, w := createTestContext(data)

	// act
	controller.Create(c)

	// assert
	assert.Equal(t, 403, w.Code)
}

func Test_Collection_ShouldFetchAll(t *testing.T) {
	// arrange
	service := createMockCollectionService()
	controller := createCollectionController(service)
	service.On("GetAll", mock.Anything).Return([]*dto.GetCollection{})
	c, w := createTestContext(nil)

	// act
	controller.All(c)

	// assert
	assert.Equal(t, 200, w.Code)
}

func Test_Collection_ShouldFetchById(t *testing.T) {
	// arrange
	service := createMockCollectionService()
	controller := createCollectionController(service)
	service.On("GetById", mock.Anything, mock.Anything).Return(&dto.GetCollection{}, nil)
	c, w := createTestContext(nil)
	c.AddParam("id", "12")

	// act
	controller.ById(c)

	// assert
	assert.Equal(t, 200, w.Code)
}

func Test_Collection_ShouldNotFetchById(t *testing.T) {
	// arrange
	service := createMockCollectionService()
	controller := createCollectionController(service)
	service.On("GetById", mock.Anything, mock.Anything).Return(nil, errors.New(""))
	c, w := createTestContext(nil)
	c.AddParam("id", "12")

	// act
	controller.ById(c)

	// assert
	assert.Equal(t, 404, w.Code)
}

func Test_Collection_ShouldAddCard(t *testing.T) {
	// arrange
	service := createMockCollectionService()
	controller := createCollectionController(service)
	service.On("AddCard", mock.Anything, mock.Anything, mock.Anything).Return(&dto.GetCollection{}, nil)
	c, w := createTestContext(&dto.CreateCardSlot{
		CardId: 0,
		Amount: 0,
	})
	c.AddParam("collectionId", "12")

	// act
	controller.AddCard(c)

	// assert
	assert.Equal(t, 200, w.Code)
}

func Test_Collection_ShouldNotAddCard(t *testing.T) {
	// arrange
	service := createMockCollectionService()
	controller := createCollectionController(service)
	service.On("AddCard", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New(""))
	c, w := createTestContext(&dto.CreateCardSlot{
		CardId: 0,
		Amount: 0,
	})
	c.AddParam("collectionId", "12")

	// act
	controller.AddCard(c)

	// assert
	assert.Equal(t, 400, w.Code)
}
