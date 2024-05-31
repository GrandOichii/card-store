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

func newCollectionController(collectionService service.CollectionService) *controller.CollectionController {
	return controller.NewCollectionController(
		collectionService,
		func(*gin.Context) {},
		func(s string, ctx *gin.Context) (string, error) {
			return "1", nil
		},
		// auth.NewJwtMiddleware(&config.Configuration{
		// 	AuthKey: "test secret key",
		// }, authService, repo).Middle.LoginHandler,
	)
}

func Test_Collection_ShouldCreate(t *testing.T) {
	// arrange
	service := newMockCollectionService()
	controller := newCollectionController(service)
	service.On("Create", mock.Anything, mock.Anything).Return(&dto.GetCollection{}, nil)
	data := dto.PostCollection{
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
	service := newMockCollectionService()
	controller := newCollectionController(service)
	service.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New(""))
	data := dto.PostCollection{
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
	s := newMockCollectionService()
	controller := newCollectionController(s)
	s.On("Create", mock.Anything, mock.Anything).Return(nil, service.ErrNotVerified)
	data := dto.PostCollection{
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
	service := newMockCollectionService()
	controller := newCollectionController(service)
	service.On("GetAll", mock.Anything).Return([]*dto.GetCollection{})
	c, w := createTestContext(nil)

	// act
	controller.All(c)

	// assert
	assert.Equal(t, 200, w.Code)
}

func Test_Collection_ShouldFetchById(t *testing.T) {
	// arrange
	service := newMockCollectionService()
	controller := newCollectionController(service)
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
	s := newMockCollectionService()
	controller := newCollectionController(s)
	s.On("GetById", mock.Anything, mock.Anything).Return(nil, service.ErrCollectionNotFound)
	c, w := createTestContext(nil)
	c.AddParam("id", "12")

	// act
	controller.ById(c)

	// assert
	assert.Equal(t, 404, w.Code)
}

func Test_Collection_ShouldEditSlot(t *testing.T) {
	// arrange
	service := newMockCollectionService()
	controller := newCollectionController(service)
	service.On("EditSlot", mock.Anything, mock.Anything, mock.Anything).Return(&dto.GetCollection{}, nil)
	c, w := createTestContext(&dto.PostCollectionSlot{
		CardId: 0,
		Amount: 0,
	})
	c.AddParam("collectionId", "12")

	// act
	controller.EditSlot(c)

	// assert
	assert.Equal(t, 200, w.Code)
}

func Test_Collection_ShouldNotEditSlot(t *testing.T) {
	// arrange
	service := newMockCollectionService()
	controller := newCollectionController(service)
	service.On("EditSlot", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New(""))
	c, w := createTestContext(&dto.PostCollectionSlot{
		CardId: 0,
		Amount: 0,
	})
	c.AddParam("collectionId", "12")

	// act
	controller.EditSlot(c)

	// assert
	assert.Equal(t, 400, w.Code)
}

func Test_Collection_ShouldNotEditSlotUnverified(t *testing.T) {
	// arrange
	s := newMockCollectionService()
	controller := newCollectionController(s)
	s.On("EditSlot", mock.Anything, mock.Anything, mock.Anything).Return(nil, service.ErrNotVerified)
	c, w := createTestContext(&dto.PostCollectionSlot{
		CardId: 0,
		Amount: 0,
	})
	c.AddParam("collectionId", "12")

	// act
	controller.EditSlot(c)

	// assert
	assert.Equal(t, 403, w.Code)
}

func Test_Collection_ShouldDelete(t *testing.T) {
	// arrange
	service := newMockCollectionService()
	controller := newCollectionController(service)
	service.On("Delete", mock.Anything, mock.Anything).Return(nil)
	c, w := createTestContext(nil)
	c.AddParam("id", "12")

	// act
	controller.Delete(c)

	// assert
	assert.Equal(t, 200, w.Code)
}

func Test_Collection_ShouldNotDelete(t *testing.T) {
	// arrange
	service := newMockCollectionService()
	controller := newCollectionController(service)
	service.On("Delete", mock.Anything, mock.Anything).Return(errors.New(""))
	c, w := createTestContext(nil)
	c.AddParam("id", "12")

	// act
	controller.Delete(c)

	// assert
	assert.Equal(t, 404, w.Code)
}

func Test_Collection_ShouldUpdateInfo(t *testing.T) {
	// arrange
	service := newMockCollectionService()
	controller := newCollectionController(service)
	service.On("UpdateInfo", mock.Anything, mock.Anything, mock.Anything).Return(&dto.GetCollection{}, nil)
	data := dto.PostCollection{
		Name:        "collection1",
		Description: "collection description",
	}
	c, w := createTestContext(data)
	c.AddParam("id", "12")

	// act
	controller.UpdateInfo(c)

	// assert
	assert.Equal(t, 200, w.Code)
}

func Test_Collection_ShouldNotUpdateInfoCollectionNotFound(t *testing.T) {
	// arrange
	s := newMockCollectionService()
	controller := newCollectionController(s)
	s.On("UpdateInfo", mock.Anything, mock.Anything, mock.Anything).Return(nil, service.ErrCollectionNotFound)
	data := dto.PostCollection{
		Name:        "collection1",
		Description: "collection description",
	}
	c, w := createTestContext(data)
	c.AddParam("id", "12")

	// act
	controller.UpdateInfo(c)

	// assert
	assert.Equal(t, 404, w.Code)
}

func Test_Collection_ShouldNotUpdateInfoBadRequest(t *testing.T) {
	// arrange
	s := newMockCollectionService()
	controller := newCollectionController(s)
	s.On("UpdateInfo", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New(""))
	data := dto.PostCollection{
		Name:        "collection1",
		Description: "collection description",
	}
	c, w := createTestContext(data)
	c.AddParam("id", "12")

	// act
	controller.UpdateInfo(c)

	// assert
	assert.Equal(t, 400, w.Code)
}
