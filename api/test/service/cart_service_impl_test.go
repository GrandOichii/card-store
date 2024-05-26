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

func newCartService(cartRepo *MockCartRepository, userRepo *MockUserRepository, cardRepo *MockCardRepository) service.CartService {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return service.NewCartServiceImpl(
		cartRepo,
		userRepo,
		cardRepo,
		validate,
	)
}

func Test_Cart_ShouldGet(t *testing.T) {
	// arrange
	cartRepo := newMockCartRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCartService(cartRepo, userRepo, cardRepo)

	userRepo.On("FindById", mock.Anything).Return(&model.User{})
	cartRepo.On("FindSingleByUserId", mock.Anything).Return(&model.Cart{})

	// act
	col, err := service.Get(1)

	// assert
	assert.NotNil(t, col)
	assert.Nil(t, err)
}

func Test_Cart_ShouldNotGet(t *testing.T) {
	// arrange
	cartRepo := newMockCartRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCartService(cartRepo, userRepo, cardRepo)

	userRepo.On("FindById", mock.Anything).Return(nil)
	cartRepo.On("FindSingleByUserId", mock.Anything).Return(&model.Cart{})

	// act
	col, err := service.Get(1)

	// assert
	assert.Nil(t, col)
	assert.NotNil(t, err)
}

func Test_Cart_ShouldEditCard(t *testing.T) {
	// arrange
	cartRepo := newMockCartRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCartService(cartRepo, userRepo, cardRepo)

	const userId uint = 1
	const cardId uint = 2

	userRepo.On("FindById", mock.Anything).Return(&model.User{})
	cardRepo.On("FindById", mock.Anything).Return(&model.Card{})
	cartRepo.On("FindSingleByUserId", mock.Anything).Return(&model.Cart{})
	cartRepo.On("Update", mock.Anything).Return(nil)

	// act
	col, err := service.EditCard(userId, &dto.PostCartSlot{
		CardId: cardId,
		Amount: 1,
	})

	// assert
	assert.NotNil(t, col)
	assert.Nil(t, err)
}

func Test_Cart_ShouldNotEditCardNegativeAmount(t *testing.T) {
	// arrange
	cartRepo := newMockCartRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCartService(cartRepo, userRepo, cardRepo)

	const userId uint = 1
	const cardId uint = 2

	userRepo.On("FindById", mock.Anything).Return(&model.User{})
	cardRepo.On("FindById", mock.Anything).Return(&model.Card{})
	cartRepo.On("FindSingleByUserId", mock.Anything).Return(&model.Cart{})
	cartRepo.On("Update", mock.Anything).Return(nil)

	// act
	col, err := service.EditCard(userId, &dto.PostCartSlot{
		CardId: cardId,
		Amount: -1,
	})

	// assert
	assert.Nil(t, col)
	assert.NotNil(t, err)
}

func Test_Cart_ShouldNotEditCardUserNotFound(t *testing.T) {
	// arrange
	cartRepo := newMockCartRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCartService(cartRepo, userRepo, cardRepo)

	const userId uint = 1
	const cardId uint = 2

	userRepo.On("FindById", mock.Anything).Return(nil)
	cardRepo.On("FindById", mock.Anything).Return(&model.Card{})
	cartRepo.On("FindSingleByUserId", mock.Anything).Return(&model.Cart{})
	cartRepo.On("Update", mock.Anything).Return(nil)

	// act
	col, err := service.EditCard(userId, &dto.PostCartSlot{
		CardId: cardId,
		Amount: 1,
	})

	// assert
	assert.Nil(t, col)
	assert.NotNil(t, err)
}

func Test_Cart_ShouldNotEditCardCardNotFound(t *testing.T) {
	// arrange
	cartRepo := newMockCartRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCartService(cartRepo, userRepo, cardRepo)

	const userId uint = 1
	const cardId uint = 2

	userRepo.On("FindById", mock.Anything).Return(&model.User{})
	cardRepo.On("FindById", mock.Anything).Return(nil)
	cartRepo.On("FindSingleByUserId", mock.Anything).Return(&model.Cart{})
	cartRepo.On("Update", mock.Anything).Return(nil)

	// act
	col, err := service.EditCard(userId, &dto.PostCartSlot{
		CardId: cardId,
		Amount: 1,
	})

	// assert
	assert.Nil(t, col)
	assert.NotNil(t, err)
}

func Test_Cart_ShouldNotEditCardFailedUpdate(t *testing.T) {
	// arrange
	cartRepo := newMockCartRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCartService(cartRepo, userRepo, cardRepo)

	const userId uint = 1
	const cardId uint = 2

	userRepo.On("FindById", mock.Anything).Return(&model.User{})
	cardRepo.On("FindById", mock.Anything).Return(&model.Card{})
	cartRepo.On("FindSingleByUserId", mock.Anything).Return(&model.Cart{})
	cartRepo.On("Update", mock.Anything).Return(errors.New(""))

	// act
	col, err := service.EditCard(userId, &dto.PostCartSlot{
		CardId: cardId,
		Amount: 1,
	})

	// assert
	assert.Nil(t, col)
	assert.NotNil(t, err)
}

// TODO? add tests for handling UpdateSlot and DeleteSlot methods returning errors

func Test_Cart_ShouldEditCardAddToAmount(t *testing.T) {
	// arrange
	cartRepo := newMockCartRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCartService(cartRepo, userRepo, cardRepo)

	const userId uint = 1
	const cardId uint = 2

	userRepo.On("FindById", mock.Anything).Return(&model.User{})
	cardRepo.On("FindById", mock.Anything).Return(&model.Card{})
	cartRepo.On("FindSingleByUserId", mock.Anything).Return(&model.Cart{
		Cards: []model.CartSlot{
			{
				CardID: cardId,
				Amount: 2,
			},
		},
	})
	cartRepo.On("Update", mock.Anything).Return(nil)
	cartRepo.On("UpdateSlot", mock.Anything).Return(nil)

	// act
	col, err := service.EditCard(userId, &dto.PostCartSlot{
		CardId: cardId,
		Amount: 1,
	})

	// assert
	assert.NotNil(t, col)
	assert.Nil(t, err)
}

func Test_Cart_ShouldEditCardSubtractFromAmount(t *testing.T) {
	// arrange
	cartRepo := newMockCartRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCartService(cartRepo, userRepo, cardRepo)

	const userId uint = 1
	const cardId uint = 2

	userRepo.On("FindById", mock.Anything).Return(&model.User{})
	cardRepo.On("FindById", mock.Anything).Return(&model.Card{})
	cartRepo.On("FindSingleByUserId", mock.Anything).Return(&model.Cart{
		Cards: []model.CartSlot{
			{
				CardID: cardId,
				Amount: 2,
			},
		},
	})
	cartRepo.On("Update", mock.Anything).Return(nil)
	cartRepo.On("UpdateSlot", mock.Anything).Return(nil)

	// act
	col, err := service.EditCard(userId, &dto.PostCartSlot{
		CardId: cardId,
		Amount: -1,
	})

	// assert
	assert.NotNil(t, col)
	assert.Nil(t, err)
}

func Test_Cart_ShouldNotEditCardSubtractFromNonexistantAmount(t *testing.T) {
	// arrange
	cartRepo := newMockCartRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCartService(cartRepo, userRepo, cardRepo)

	const userId uint = 1
	const cardId uint = 2

	userRepo.On("FindById", mock.Anything).Return(&model.User{})
	cardRepo.On("FindById", mock.Anything).Return(&model.Card{})
	cartRepo.On("FindSingleByUserId", mock.Anything).Return(&model.Cart{
		Cards: []model.CartSlot{
			{
				CardID: cardId,
				Amount: 2,
			},
		},
	})
	cartRepo.On("Update", mock.Anything).Return(nil)
	cartRepo.On("UpdateSlot", mock.Anything).Return(nil)

	// act
	col, err := service.EditCard(userId, &dto.PostCartSlot{
		CardId: 3,
		Amount: -1,
	})

	// assert
	assert.Nil(t, col)
	assert.NotNil(t, err)
}

func Test_Cart_ShouldEditCardDeleteSlot(t *testing.T) {
	// arrange
	cartRepo := newMockCartRepository()
	userRepo := newMockUserRepository()
	cardRepo := newMockCardRepository()
	service := newCartService(cartRepo, userRepo, cardRepo)

	const userId uint = 1
	const cardId uint = 2

	userRepo.On("FindById", mock.Anything).Return(&model.User{})
	cardRepo.On("FindById", mock.Anything).Return(&model.Card{})
	cartRepo.On("FindSingleByUserId", mock.Anything).Return(&model.Cart{
		Cards: []model.CartSlot{
			{
				CardID: cardId,
				Amount: 2,
			},
		},
	})
	cartRepo.On("Update", mock.Anything).Return(nil)
	cartRepo.On("DeleteSlot", mock.Anything).Return(nil)

	// act
	col, err := service.EditCard(userId, &dto.PostCartSlot{
		CardId: cardId,
		Amount: -2,
	})

	// assert
	assert.NotNil(t, col)
	assert.Nil(t, err)
}
