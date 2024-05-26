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

func newCollectionService(collectionRepo *MockCollectionRepository, userRepo *MockUserRepository, cardRepo *MockCardRepository) service.CollectionService {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return service.NewCollectionServiceImpl(
		collectionRepo,
		userRepo,
		cardRepo,
		validate,
	)
}

func Test_Collection_ShouldCreate(t *testing.T) {
	// arrange
	colRepo := newMockCollectionRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCollectionService(colRepo, userRepo, cardRepo)

	userRepo.On("FindById", mock.Anything).Return(&model.User{Verified: true})
	colRepo.On("Save", mock.Anything).Return(nil)

	// act
	col, err := service.Create(&dto.PostCollection{
		Name:        "collection1",
		Description: "collection description",
	}, 1)

	// assert
	assert.NotNil(t, col)
	assert.Nil(t, err)
}

func Test_Collection_ShouldNotCreateUnverified(t *testing.T) {
	// arrange
	colRepo := newMockCollectionRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	s := newCollectionService(colRepo, userRepo, cardRepo)

	userRepo.On("FindById", mock.Anything).Return(&model.User{Verified: false})
	colRepo.On("Save", mock.Anything).Return(nil)

	// act
	col, err := s.Create(&dto.PostCollection{
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
	colRepo := newMockCollectionRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCollectionService(colRepo, userRepo, cardRepo)

	userRepo.On("FindById", mock.Anything).Return(nil)
	colRepo.On("Save", mock.Anything).Return(nil)

	// act
	col, err := service.Create(&dto.PostCollection{
		Name:        "collection1",
		Description: "collection description",
	}, 1)

	// assert
	assert.Nil(t, col)
	assert.NotNil(t, err)
}

func Test_Collection_ShouldNotCreateSaveFail(t *testing.T) {
	// arrange
	colRepo := newMockCollectionRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCollectionService(colRepo, userRepo, cardRepo)

	userRepo.On("FindById", mock.Anything).Return(&model.User{Verified: true})
	colRepo.On("Save", mock.Anything).Return(errors.New(""))

	// act
	col, err := service.Create(&dto.PostCollection{
		Name:        "collection1",
		Description: "collection description",
	}, 1)

	// assert
	assert.Nil(t, col)
	assert.NotNil(t, err)
}

func Test_Collection_ShouldGetAll(t *testing.T) {
	// arrange
	colRepo := newMockCollectionRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCollectionService(colRepo, userRepo, cardRepo)

	colRepo.On("FindByOwnerId", mock.Anything).Return([]*model.Collection{})

	// act
	cards := service.GetAll(1)

	// assert
	assert.NotNil(t, cards)
}

func Test_Collection_ShouldGetById(t *testing.T) {
	// arrange
	colRepo := newMockCollectionRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCollectionService(colRepo, userRepo, cardRepo)
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
	colRepo := newMockCollectionRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCollectionService(colRepo, userRepo, cardRepo)
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
	colRepo := newMockCollectionRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCollectionService(colRepo, userRepo, cardRepo)
	const userId uint = 1

	colRepo.On("FindById", mock.Anything).Return(&model.Collection{OwnerID: 2})

	// act
	col, err := service.GetById(2, userId)

	// assert
	assert.Nil(t, col)
	assert.NotNil(t, err)
}

func Test_Collection_ShouldEditSlot(t *testing.T) {
	// arrange
	colRepo := newMockCollectionRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCollectionService(colRepo, userRepo, cardRepo)

	const colId uint = 1
	const userId uint = 2

	userRepo.On("FindById", mock.Anything).Return(&model.User{Verified: true})
	colRepo.On("FindById", mock.Anything).Return(&model.Collection{OwnerID: userId})
	colRepo.On("Update", mock.Anything).Return(nil)
	cardRepo.On("FindById", mock.Anything).Return(&model.Card{})

	// act
	col, err := service.EditSlot(&dto.PostCollectionSlot{
		CardId: 1,
		Amount: 1,
	}, colId, userId)

	// assert
	assert.NotNil(t, col)
	assert.Nil(t, err)
}

func Test_Collection_ShouldNotEditSlotNotFound(t *testing.T) {
	// arrange
	colRepo := newMockCollectionRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCollectionService(colRepo, userRepo, cardRepo)

	const colId uint = 1
	const userId uint = 2

	userRepo.On("FindById", mock.Anything).Return(&model.User{Verified: true})
	colRepo.On("FindById", mock.Anything).Return(&model.Collection{OwnerID: userId})
	colRepo.On("Update", mock.Anything).Return(nil)
	cardRepo.On("FindById", mock.Anything).Return(nil)

	// act
	col, err := service.EditSlot(&dto.PostCollectionSlot{
		CardId: 1,
		Amount: 1,
	}, colId, userId)

	// assert
	assert.Nil(t, col)
	assert.NotNil(t, err)
}

func Test_Collection_ShouldNotEditSlotUnverified(t *testing.T) {
	// arrange
	colRepo := newMockCollectionRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	s := newCollectionService(colRepo, userRepo, cardRepo)

	const colId uint = 1
	const userId uint = 2

	userRepo.On("FindById", mock.Anything).Return(&model.User{Verified: false})
	colRepo.On("FindById", mock.Anything).Return(&model.Collection{OwnerID: userId})
	colRepo.On("Update", mock.Anything).Return(nil)
	cardRepo.On("FindById", mock.Anything).Return(&model.Card{})

	// act
	col, err := s.EditSlot(&dto.PostCollectionSlot{
		CardId: 1,
		Amount: 1,
	}, colId, userId)

	// assert
	assert.Nil(t, col)
	assert.NotNil(t, err)
	assert.Equal(t, service.ErrNotVerified, err)
}

func Test_Collection_ShouldNotEditSlotNoCollection(t *testing.T) {
	// arrange
	colRepo := newMockCollectionRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCollectionService(colRepo, userRepo, cardRepo)

	const colId uint = 1
	const userId uint = 2

	userRepo.On("FindById", mock.Anything).Return(&model.User{Verified: true})
	colRepo.On("FindById", mock.Anything).Return(nil)
	colRepo.On("Update", mock.Anything).Return(nil)
	cardRepo.On("FindById", mock.Anything).Return(&model.Card{})

	// act
	col, err := service.EditSlot(&dto.PostCollectionSlot{
		CardId: 1,
		Amount: 1,
	}, colId, userId)

	// assert
	assert.Nil(t, col)
	assert.NotNil(t, err)
}

func Test_Collection_ShouldNotEditSlotNoUser(t *testing.T) {
	// arrange
	colRepo := newMockCollectionRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCollectionService(colRepo, userRepo, cardRepo)

	const colId uint = 1
	const userId uint = 2

	userRepo.On("FindById", mock.Anything).Return(&model.User{Verified: true})
	colRepo.On("FindById", mock.Anything).Return(nil)
	colRepo.On("Update", mock.Anything).Return(nil)
	cardRepo.On("FindById", mock.Anything).Return(&model.Card{})

	// act
	col, err := service.EditSlot(&dto.PostCollectionSlot{
		CardId: 1,
		Amount: 1,
	}, colId, userId)

	// assert
	assert.Nil(t, col)
	assert.NotNil(t, err)
}

func Test_Collection_ShouldNotEditSlotMismathOwnerId(t *testing.T) {
	// arrange
	colRepo := newMockCollectionRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCollectionService(colRepo, userRepo, cardRepo)

	const colId uint = 1
	const userId uint = 2

	userRepo.On("FindById", mock.Anything).Return(&model.User{Verified: true})
	colRepo.On("FindById", mock.Anything).Return(&model.Collection{})
	colRepo.On("Update", mock.Anything).Return(nil)
	cardRepo.On("FindById", mock.Anything).Return(&model.Card{})

	// act
	col, err := service.EditSlot(&dto.PostCollectionSlot{
		CardId: 1,
		Amount: 1,
	}, colId, userId)

	// assert
	assert.Nil(t, col)
	assert.NotNil(t, err)
}

func Test_Collection_ShouldNotEditSlotFailedUpdate(t *testing.T) {
	// arrange
	colRepo := newMockCollectionRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCollectionService(colRepo, userRepo, cardRepo)

	const colId uint = 1
	const userId uint = 2

	userRepo.On("FindById", mock.Anything).Return(&model.User{Verified: true})
	colRepo.On("FindById", mock.Anything).Return(&model.Collection{OwnerID: userId})
	colRepo.On("Update", mock.Anything).Return(errors.New(""))
	cardRepo.On("FindById", mock.Anything).Return(&model.Card{})

	// act
	col, err := service.EditSlot(&dto.PostCollectionSlot{
		CardId: 1,
		Amount: 1,
	}, colId, userId)

	// assert
	assert.Nil(t, col)
	assert.NotNil(t, err)
}

func Test_Collection_ShouldNotEditSlotAmountZero(t *testing.T) {
	// arrange
	colRepo := newMockCollectionRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCollectionService(colRepo, userRepo, cardRepo)

	const colId uint = 1
	const userId uint = 2

	userRepo.On("FindById", mock.Anything).Return(&model.User{Verified: true})
	colRepo.On("FindById", mock.Anything).Return(&model.Collection{OwnerID: userId})
	colRepo.On("Update", mock.Anything).Return(nil)
	cardRepo.On("FindById", mock.Anything).Return(&model.Card{})

	// act
	col, err := service.EditSlot(&dto.PostCollectionSlot{
		CardId: 1,
		Amount: 0,
	}, colId, userId)

	// assert
	assert.Nil(t, col)
	assert.NotNil(t, err)
}

// TODO? add tests for handling UpdateSlot and DeleteSlot methods returning errors

func Test_Collection_ShouldEditSlotAddToAmount(t *testing.T) {
	// arrange
	colRepo := newMockCollectionRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCollectionService(colRepo, userRepo, cardRepo)

	const colId uint = 1
	const userId uint = 2
	const cardId uint = 3

	userRepo.On("FindById", mock.Anything).Return(&model.User{Verified: true})
	colRepo.On("FindById", mock.Anything).Return(&model.Collection{
		OwnerID: userId,
		Cards: []model.CollectionSlot{
			{
				CardID: cardId,
				Amount: 10,
			},
		},
	})
	colRepo.On("Update", mock.Anything).Return(nil)
	colRepo.On("UpdateSlot", mock.Anything).Return(nil)
	cardRepo.On("FindById", mock.Anything).Return(&model.Card{})

	// act
	col, err := service.EditSlot(&dto.PostCollectionSlot{
		CardId: cardId,
		Amount: 3,
	}, colId, userId)

	// assert
	assert.NotNil(t, col)
	assert.Nil(t, err)
}

func Test_Collection_ShouldEditSlotSubtractFromAmount(t *testing.T) {
	// arrange
	colRepo := newMockCollectionRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCollectionService(colRepo, userRepo, cardRepo)

	const colId uint = 1
	const userId uint = 2
	const cardId uint = 3

	userRepo.On("FindById", mock.Anything).Return(&model.User{Verified: true})
	colRepo.On("FindById", mock.Anything).Return(&model.Collection{
		OwnerID: userId,
		Cards: []model.CollectionSlot{
			{
				CardID: cardId,
				Amount: 10,
			},
		},
	})
	colRepo.On("Update", mock.Anything).Return(nil)
	colRepo.On("UpdateSlot", mock.Anything).Return(nil)
	cardRepo.On("FindById", mock.Anything).Return(&model.Card{})

	// act
	col, err := service.EditSlot(&dto.PostCollectionSlot{
		CardId: cardId,
		Amount: -3,
	}, colId, userId)

	// assert
	assert.NotNil(t, col)
	assert.Nil(t, err)
}

func Test_Collection_ShouldNotEditSlotSubtractFromNonexistantAmount(t *testing.T) {
	// arrange
	colRepo := newMockCollectionRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCollectionService(colRepo, userRepo, cardRepo)

	const colId uint = 1
	const userId uint = 2
	const cardId uint = 3

	userRepo.On("FindById", mock.Anything).Return(&model.User{Verified: true})
	colRepo.On("FindById", mock.Anything).Return(&model.Collection{
		OwnerID: userId,
	})
	colRepo.On("Update", mock.Anything).Return(nil)
	colRepo.On("UpdateSlot", mock.Anything).Return(nil)
	cardRepo.On("FindById", mock.Anything).Return(&model.Card{})

	// act
	col, err := service.EditSlot(&dto.PostCollectionSlot{
		CardId: cardId,
		Amount: -3,
	}, colId, userId)

	// assert
	assert.Nil(t, col)
	assert.NotNil(t, err)
}

func Test_Collection_ShouldEditSlotDeleteSlot(t *testing.T) {
	// arrange
	colRepo := newMockCollectionRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCollectionService(colRepo, userRepo, cardRepo)

	const colId uint = 1
	const userId uint = 2
	const cardId uint = 3

	userRepo.On("FindById", mock.Anything).Return(&model.User{Verified: true})
	colRepo.On("FindById", mock.Anything).Return(&model.Collection{
		OwnerID: userId,
		Cards: []model.CollectionSlot{
			{
				CardID: cardId,
				Amount: 10,
			},
		},
	})
	colRepo.On("Update", mock.Anything).Return(nil)
	colRepo.On("UpdateSlot", mock.Anything).Return(nil)
	colRepo.On("DeleteSlot", mock.Anything).Return(nil)
	cardRepo.On("FindById", mock.Anything).Return(&model.Card{})

	// act
	col, err := service.EditSlot(&dto.PostCollectionSlot{
		CardId: cardId,
		Amount: -10,
	}, colId, userId)

	// assert
	assert.NotNil(t, col)
	assert.Nil(t, err)
}

func Test_Collection_ShouldDelete(t *testing.T) {
	// arrange
	colRepo := newMockCollectionRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCollectionService(colRepo, userRepo, cardRepo)
	const userId uint = 1

	colRepo.On("FindById", mock.Anything).Return(&model.Collection{OwnerID: userId})
	colRepo.On("Delete", mock.Anything).Return(nil)

	// act
	err := service.Delete(2, userId)

	// assert
	assert.Nil(t, err)
}

func Test_Collection_ShouldNotDeleteMismatchUserId(t *testing.T) {
	// arrange
	colRepo := newMockCollectionRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCollectionService(colRepo, userRepo, cardRepo)
	const userId uint = 1

	colRepo.On("FindById", mock.Anything).Return(&model.Collection{OwnerID: 2})
	colRepo.On("Delete", mock.Anything).Return(nil)

	// act
	err := service.Delete(2, userId)

	// assert
	assert.NotNil(t, err)
}

func Test_Collection_ShouldNotDeleteNotFound(t *testing.T) {
	// arrange
	colRepo := newMockCollectionRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCollectionService(colRepo, userRepo, cardRepo)
	const userId uint = 1

	colRepo.On("FindById", mock.Anything).Return(&model.Collection{OwnerID: userId})
	colRepo.On("Delete", mock.Anything).Return(errors.New(""))

	// act
	err := service.Delete(2, userId)

	// assert
	assert.NotNil(t, err)
}

func Test_Collection_ShouldUpdateInfo(t *testing.T) {
	// arrange
	colRepo := newMockCollectionRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCollectionService(colRepo, userRepo, cardRepo)
	const userId uint = 1

	colRepo.On("FindById", mock.Anything).Return(&model.Collection{OwnerID: userId})
	colRepo.On("Update", mock.Anything).Return(nil)
	userRepo.On("FindById", mock.Anything).Return(&model.User{Verified: true})

	// act
	col, err := service.UpdateInfo(&dto.PostCollection{
		Name:        "collection1",
		Description: "collection description",
	}, 2, userId)

	// assert
	assert.NotNil(t, col)
	assert.Nil(t, err)
}

func Test_Collection_ShouldNotUpdateInfoNotVerified(t *testing.T) {
	// arrange
	colRepo := newMockCollectionRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	s := newCollectionService(colRepo, userRepo, cardRepo)
	const userId uint = 1

	colRepo.On("FindById", mock.Anything).Return(&model.Collection{OwnerID: userId})
	colRepo.On("Update", mock.Anything).Return(nil)
	userRepo.On("FindById", mock.Anything).Return(&model.User{Verified: false})

	// act
	col, err := s.UpdateInfo(&dto.PostCollection{
		Name:        "collection1",
		Description: "collection description",
	}, 2, userId)

	// assert
	assert.Nil(t, col)
	assert.Equal(t, service.ErrNotVerified, err)
}

func Test_Collection_ShouldNotUpdateInfoNoCollection(t *testing.T) {
	// arrange
	colRepo := newMockCollectionRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	s := newCollectionService(colRepo, userRepo, cardRepo)
	const userId uint = 1

	colRepo.On("Update", mock.Anything).Return(nil)
	userRepo.On("FindById", mock.Anything).Return(&model.User{Verified: true})
	colRepo.On("FindById", mock.Anything).Return(nil)

	// act
	col, err := s.UpdateInfo(&dto.PostCollection{
		Name:        "collection1",
		Description: "collection description",
	}, 2, userId)

	// assert
	assert.Nil(t, col)
	assert.Equal(t, service.ErrCollectionNotFound, err)
}

func Test_Collection_ShouldNotUpdateInfoBadData(t *testing.T) {
	// arrange
	colRepo := newMockCollectionRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCollectionService(colRepo, userRepo, cardRepo)
	const userId uint = 1

	colRepo.On("FindById", mock.Anything).Return(&model.Collection{OwnerID: userId})
	colRepo.On("Update", mock.Anything).Return(nil)
	userRepo.On("FindById", mock.Anything).Return(&model.User{Verified: true})

	// act
	col, err := service.UpdateInfo(&dto.PostCollection{
		Name:        "",
		Description: "collection description",
	}, 2, userId)

	// assert
	assert.Nil(t, col)
	assert.NotNil(t, err)
}

func Test_Collection_ShouldNotUpdateInfoBadUpdate(t *testing.T) {
	// arrange
	colRepo := newMockCollectionRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCollectionService(colRepo, userRepo, cardRepo)
	const userId uint = 1

	colRepo.On("FindById", mock.Anything).Return(&model.Collection{OwnerID: userId})
	colRepo.On("Update", mock.Anything).Return(errors.New(""))
	userRepo.On("FindById", mock.Anything).Return(&model.User{Verified: true})

	// act
	col, err := service.UpdateInfo(&dto.PostCollection{
		Name:        "collection1",
		Description: "collection description",
	}, 2, userId)

	// assert
	assert.Nil(t, col)
	assert.NotNil(t, err)
}
