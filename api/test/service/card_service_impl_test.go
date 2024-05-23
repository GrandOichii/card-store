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

func createCardService(cardRepo *MockCardRepository, userRepo *MockUserRepository) service.CardService {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return service.NewCardServiceImpl(
		cardRepo,
		userRepo,
		validate,
	)
}

func Test_ShouldAdd(t *testing.T) {
	// arrange
	cardRepo := createMockCardRepository()
	userRepo := createMockUserRepository()
	service := createCardService(cardRepo, userRepo)

	cardRepo.On("Save", mock.Anything).Return(nil)
	userRepo.On("FindByUsername", mock.Anything).Return(&model.User{})

	// act
	card, err := service.Add(&dto.CreateCard{
		Name:  "card name",
		Text:  "card text",
		Price: 10,
		Type:  "CT1",
	}, "userID")

	// assert
	assert.NotNil(t, card)
	assert.Nil(t, err)
}

func Test_ShouldNotAdd(t *testing.T) {
	// arrange
	cardRepo := createMockCardRepository()
	userRepo := createMockUserRepository()
	service := createCardService(cardRepo, userRepo)

	cardRepo.On("Save", mock.Anything).Return(errors.New(""))
	userRepo.On("FindByUsername", mock.Anything).Return(&model.User{})

	// act
	card, err := service.Add(&dto.CreateCard{
		Name:  "card name",
		Text:  "card text",
		Price: 10,
		Type:  "CT1",
	}, "userID")

	// assert
	assert.Nil(t, card)
	assert.NotNil(t, err)

}
func Test_ShouldNotAddUnknownUser(t *testing.T) {
	// arrange
	cardRepo := createMockCardRepository()
	userRepo := createMockUserRepository()
	service := createCardService(cardRepo, userRepo)

	cardRepo.On("Save", mock.Anything).Return(nil)
	userRepo.On("FindByUsername", mock.Anything).Return(nil)

	// act
	card, err := service.Add(&dto.CreateCard{
		Name:  "card name",
		Text:  "card text",
		Price: 10,
		Type:  "CT1",
	}, "userID")

	// assert
	assert.Nil(t, card)
	assert.NotNil(t, err)
}

func Test_ShouldGetById(t *testing.T) {
	// arrange
	cardRepo := createMockCardRepository()
	userRepo := createMockUserRepository()
	service := createCardService(cardRepo, userRepo)

	cardRepo.On("FindById", mock.Anything).Return(&model.Card{})

	// act
	card, err := service.GetById(1)

	// assert
	assert.NotNil(t, card)
	assert.Nil(t, err)
}

func Test_ShouldNotGetById(t *testing.T) {
	// arrange
	cardRepo := createMockCardRepository()
	userRepo := createMockUserRepository()
	service := createCardService(cardRepo, userRepo)

	cardRepo.On("FindById", mock.Anything).Return(nil)

	// act
	card, err := service.GetById(1)

	// assert
	assert.Nil(t, card)
	assert.NotNil(t, err)
}

func Test_ShouldGetByType(t *testing.T) {
	// arrange
	cardRepo := createMockCardRepository()
	userRepo := createMockUserRepository()
	service := createCardService(cardRepo, userRepo)

	cardRepo.On("FindByType", mock.Anything).Return([]*model.Card{}, nil)

	// act
	cards := service.GetByType("CT1")

	// assert
	assert.NotNil(t, cards)
}
