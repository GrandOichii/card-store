package service_test

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"store.api/dto"
	"store.api/model"
	"store.api/service"
)

func createCollectionService(collectionRepo *MockCollectionRepository, userRepo *MockUserRepository) service.CollectionService {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return service.NewCollectionServiceImpl(
		collectionRepo,
		userRepo,
		validate,
	)
}

func Test_Collection_ShouldCreate(t *testing.T) {
	// arrange
	colRepo := createMockCollectionRepository()
	userRepo := createMockUserRepository()
	service := createCollectionService(colRepo, userRepo)

	userRepo.On("FindById", mock.Anything).Return(&model.User{Verified: true})
	colRepo.On("Save", mock.Anything).Return(nil)

	// act
	col, err := service.Create(&dto.CreateCollection{
		Name:        "collection1",
		Description: "collection description",
	}, 1)

	// assert
	assert.NotNil(t, col)
	assert.Nil(t, err)
}

func Test_Collection_ShouldNotCreateUnverified(t *testing.T) {
	// arrange
	colRepo := createMockCollectionRepository()
	userRepo := createMockUserRepository()
	s := createCollectionService(colRepo, userRepo)

	userRepo.On("FindById", mock.Anything).Return(&model.User{Verified: false})
	colRepo.On("Save", mock.Anything).Return(nil)

	// act
	col, err := s.Create(&dto.CreateCollection{
		Name:        "collection1",
		Description: "collection description",
	}, 1)

	// assert
	assert.Nil(t, col)
	assert.NotNil(t, err)
	assert.Equal(t, service.ErrNotVerified, err)
}

func Test_Collection_ShouldNotCreateInvalidUser(t *testing.T) {
	// arrange
	colRepo := createMockCollectionRepository()
	userRepo := createMockUserRepository()
	service := createCollectionService(colRepo, userRepo)

	userRepo.On("FindById", mock.Anything).Return(nil)
	colRepo.On("Save", mock.Anything).Return(nil)

	// act
	col, err := service.Create(&dto.CreateCollection{
		Name:        "collection1",
		Description: "collection description",
	}, 1)

	// assert
	assert.Nil(t, col)
	assert.NotNil(t, err)
}

func Test_Collection_ShouldNotCreateSaveFail(t *testing.T) {
	// arrange
	colRepo := createMockCollectionRepository()
	userRepo := createMockUserRepository()
	service := createCollectionService(colRepo, userRepo)

	userRepo.On("FindById", mock.Anything).Return(&model.User{Verified: true})
	colRepo.On("Save", mock.Anything).Return(errors.New(""))

	// act
	col, err := service.Create(&dto.CreateCollection{
		Name:        "collection1",
		Description: "collection description",
	}, 1)

	// assert
	assert.Nil(t, col)
	assert.NotNil(t, err)
}

func Test_Collection_ShouldGetAll(t *testing.T) {
	// arrange
	colRepo := createMockCollectionRepository()
	userRepo := createMockUserRepository()
	service := createCollectionService(colRepo, userRepo)

	colRepo.On("FindByOwnerId", mock.Anything).Return([]*model.Collection{})

	// act
	cards := service.GetAll(1)

	// assert
	assert.NotNil(t, cards)
}

func Test_Collection_ShouldGetById(t *testing.T) {
	// arrange
	colRepo := createMockCollectionRepository()
	userRepo := createMockUserRepository()
	service := createCollectionService(colRepo, userRepo)
	const userId uint = 1

	colRepo.On("FindById", mock.Anything).Return(&model.Collection{OwnerID: userId})

	// act
	col, err := service.GetById(2, userId)

	// assert
	assert.NotNil(t, col)
	assert.Nil(t, err)
}

func Test_Collection_ShouldNotGetById(t *testing.T) {
	// arrange
	colRepo := createMockCollectionRepository()
	userRepo := createMockUserRepository()
	service := createCollectionService(colRepo, userRepo)
	const userId uint = 1

	colRepo.On("FindById", mock.Anything).Return(nil)

	// act
	col, err := service.GetById(2, userId)

	// assert
	assert.Nil(t, col)
	assert.NotNil(t, err)
}

func Test_Collection_ShouldNotGetByIdOwnerMismatch(t *testing.T) {
	// arrange
	colRepo := createMockCollectionRepository()
	userRepo := createMockUserRepository()
	service := createCollectionService(colRepo, userRepo)
	const userId uint = 1

	colRepo.On("FindById", mock.Anything).Return(&model.Collection{OwnerID: 2})

	// act
	col, err := service.GetById(2, userId)

	// assert
	assert.Nil(t, col)
	assert.NotNil(t, err)
}

// TODO test user verification when adding card to collection

func Test_Collection_ShouldAddCard(t *testing.T) {
	// arrange
	colRepo := createMockCollectionRepository()
	userRepo := createMockUserRepository()
	service := createCollectionService(colRepo, userRepo)

	const colId uint = 1
	const userId uint = 2

	userRepo.On("FindById", mock.Anything).Return(&model.User{Verified: true})
	colRepo.On("FindById", mock.Anything).Return(&model.Collection{OwnerID: userId})
	colRepo.On("Update", mock.Anything).Return(nil)

	// act
	col, err := service.AddCard(&dto.CreateCardSlot{
		CardId: 1,
		Amount: 1,
	}, colId, userId)

	// assert
	assert.NotNil(t, col)
	assert.Nil(t, err)
}

func Test_Collection_ShouldNotAddCardUnverified(t *testing.T) {
	// arrange
	colRepo := createMockCollectionRepository()
	userRepo := createMockUserRepository()
	s := createCollectionService(colRepo, userRepo)

	const colId uint = 1
	const userId uint = 2

	userRepo.On("FindById", mock.Anything).Return(&model.User{Verified: false})
	colRepo.On("FindById", mock.Anything).Return(&model.Collection{OwnerID: userId})
	colRepo.On("Update", mock.Anything).Return(nil)

	// act
	col, err := s.AddCard(&dto.CreateCardSlot{
		CardId: 1,
		Amount: 1,
	}, colId, userId)

	// assert
	assert.Nil(t, col)
	assert.NotNil(t, err)
	assert.Equal(t, service.ErrNotVerified, err)
}

func Test_Collection_ShouldNotAddCardNoCollection(t *testing.T) {
	// arrange
	colRepo := createMockCollectionRepository()
	userRepo := createMockUserRepository()
	service := createCollectionService(colRepo, userRepo)

	const colId uint = 1
	const userId uint = 2

	userRepo.On("FindById", mock.Anything).Return(&model.User{Verified: true})
	colRepo.On("FindById", mock.Anything).Return(nil)
	colRepo.On("Update", mock.Anything).Return(nil)

	// act
	col, err := service.AddCard(&dto.CreateCardSlot{
		CardId: 1,
		Amount: 1,
	}, colId, userId)

	// assert
	assert.Nil(t, col)
	assert.NotNil(t, err)
}

func Test_Collection_ShouldNotAddCardNoUser(t *testing.T) {
	// arrange
	colRepo := createMockCollectionRepository()
	userRepo := createMockUserRepository()
	service := createCollectionService(colRepo, userRepo)

	const colId uint = 1
	const userId uint = 2

	userRepo.On("FindById", mock.Anything).Return(&model.User{Verified: true})
	colRepo.On("FindById", mock.Anything).Return(nil)
	colRepo.On("Update", mock.Anything).Return(nil)

	// act
	col, err := service.AddCard(&dto.CreateCardSlot{
		CardId: 1,
		Amount: 1,
	}, colId, userId)

	// assert
	assert.Nil(t, col)
	assert.NotNil(t, err)
}

func Test_Collection_ShouldNotAddCardMismathOwnerId(t *testing.T) {
	// arrange
	colRepo := createMockCollectionRepository()
	userRepo := createMockUserRepository()
	service := createCollectionService(colRepo, userRepo)

	const colId uint = 1
	const userId uint = 2

	userRepo.On("FindById", mock.Anything).Return(&model.User{Verified: true})
	colRepo.On("FindById", mock.Anything).Return(&model.Collection{})
	colRepo.On("Update", mock.Anything).Return(nil)

	// act
	col, err := service.AddCard(&dto.CreateCardSlot{
		CardId: 1,
		Amount: 1,
	}, colId, userId)

	// assert
	assert.Nil(t, col)
	assert.NotNil(t, err)
}

func Test_Collection_ShouldNotAddCardFailedUpdate(t *testing.T) {
	// arrange
	colRepo := createMockCollectionRepository()
	userRepo := createMockUserRepository()
	service := createCollectionService(colRepo, userRepo)

	const colId uint = 1
	const userId uint = 2

	userRepo.On("FindById", mock.Anything).Return(&model.User{Verified: true})
	colRepo.On("FindById", mock.Anything).Return(&model.Collection{OwnerID: userId})
	colRepo.On("Update", mock.Anything).Return(errors.New(""))

	// act
	col, err := service.AddCard(&dto.CreateCardSlot{
		CardId: 1,
		Amount: 1,
	}, colId, userId)

	// assert
	assert.Nil(t, col)
	assert.NotNil(t, err)
}

func Test_Collection_ShouldAddCardAmountZero(t *testing.T) {
	// arrange
	colRepo := createMockCollectionRepository()
	userRepo := createMockUserRepository()
	service := createCollectionService(colRepo, userRepo)

	const colId uint = 1
	const userId uint = 2

	userRepo.On("FindById", mock.Anything).Return(&model.User{Verified: true})
	colRepo.On("FindById", mock.Anything).Return(&model.Collection{OwnerID: userId})
	colRepo.On("Update", mock.Anything).Return(nil)

	// act
	col, err := service.AddCard(&dto.CreateCardSlot{
		CardId: 1,
		Amount: 0,
	}, colId, userId)

	// assert
	assert.Nil(t, col)
	assert.NotNil(t, err)
}
