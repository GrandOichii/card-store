package service_test

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"store.api/dto"
	"store.api/model"
	"store.api/query"
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

func Test_Card_ShouldAdd(t *testing.T) {
	// arrange
	cardRepo := createMockCardRepository()
	userRepo := createMockUserRepository()
	service := createCardService(cardRepo, userRepo)

	cardRepo.On("Save", mock.Anything).Return(nil)
	userRepo.On("FindById", mock.Anything).Return(&model.User{})

	// act
	card, err := service.Add(&dto.CreateCard{
		Name:     "card name",
		Text:     "card text",
		Price:    10,
		Type:     "CT1",
		Language: "ENG",
	}, 1)

	// assert
	assert.NotNil(t, card)
	assert.Nil(t, err)
}

func Test_Card_ShouldNotAdd(t *testing.T) {
	// arrange
	cardRepo := createMockCardRepository()
	userRepo := createMockUserRepository()
	service := createCardService(cardRepo, userRepo)

	cardRepo.On("Save", mock.Anything).Return(errors.New(""))
	userRepo.On("FindById", mock.Anything).Return(&model.User{})

	// act
	card, err := service.Add(&dto.CreateCard{
		Name:  "card name",
		Text:  "card text",
		Price: 10,
		Type:  "CT1",
	}, 1)

	// assert
	assert.Nil(t, card)
	assert.NotNil(t, err)

}
func Test_Card_ShouldNotAddUnknownUser(t *testing.T) {
	// arrange
	cardRepo := createMockCardRepository()
	userRepo := createMockUserRepository()
	service := createCardService(cardRepo, userRepo)

	cardRepo.On("Save", mock.Anything).Return(nil)
	userRepo.On("FindById", mock.Anything).Return(nil)

	// act
	card, err := service.Add(&dto.CreateCard{
		Name:  "card name",
		Text:  "card text",
		Price: 10,
		Type:  "CT1",
	}, 1)

	// assert
	assert.Nil(t, card)
	assert.NotNil(t, err)
}

func Test_Card_ShouldGetById(t *testing.T) {
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

func Test_Card_ShouldNotGetById(t *testing.T) {
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

func Test_Card_ShouldGetByQuery(t *testing.T) {
	// arrange
	cardRepo := createMockCardRepository()
	userRepo := createMockUserRepository()
	service := createCardService(cardRepo, userRepo)

	cardRepo.On("Query", mock.Anything, mock.Anything).Return([]*model.Card{}, nil)

	// act
	cards := service.Query(&query.CardQuery{})

	// assert
	assert.NotNil(t, cards)
}
