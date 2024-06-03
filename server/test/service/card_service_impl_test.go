package service_test

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"store.api/config"
	"store.api/dto"
	"store.api/model"
	"store.api/query"
	"store.api/service"
)

func newCardService(cardRepo *MockCardRepository, userRepo *MockUserRepository, langRepo *MockLanguageRepository) service.CardService {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return service.NewCardServiceImpl(
		&config.Configuration{
			Db: config.DbConfiguration{
				Cards: config.CardsDbConfiguration{
					PageSize: 30,
				},
			},
		},
		cardRepo,
		userRepo,
		langRepo,
		validate,
	)
}

func Test_Card_ShouldAdd(t *testing.T) {
	// arrange
	cardRepo := newMockCardRepository()
	userRepo := newMockUserRepository()
	langRepo := newMockLanguageRepository()
	service := newCardService(cardRepo, userRepo, langRepo)

	cardRepo.On("Save", mock.Anything).Return(nil)
	cardRepo.On("FindById", mock.Anything).Return(&model.Card{})
	userRepo.On("FindById", mock.Anything).Return(&model.User{})

	// act
	card, err := service.Add(&dto.PostCard{
		Name:      "card name",
		Text:      "card text",
		Price:     10,
		Type:      "CT1",
		Language:  "ENG",
		Key:       "key1",
		Expansion: "exp1",
	}, 1)

	// assert
	assert.NotNil(t, card)
	assert.Nil(t, err)
}

func Test_Card_ShouldNotAdd(t *testing.T) {
	// arrange
	cardRepo := newMockCardRepository()
	userRepo := newMockUserRepository()
	langRepo := newMockLanguageRepository()
	service := newCardService(cardRepo, userRepo, langRepo)

	cardRepo.On("Save", mock.Anything).Return(errors.New(""))
	userRepo.On("FindById", mock.Anything).Return(&model.User{})

	// act
	card, err := service.Add(&dto.PostCard{
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
	cardRepo := newMockCardRepository()
	userRepo := newMockUserRepository()
	langRepo := newMockLanguageRepository()
	service := newCardService(cardRepo, userRepo, langRepo)

	cardRepo.On("Save", mock.Anything).Return(nil)
	userRepo.On("FindById", mock.Anything).Return(nil)

	// act
	card, err := service.Add(&dto.PostCard{
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
	cardRepo := newMockCardRepository()
	userRepo := newMockUserRepository()
	langRepo := newMockLanguageRepository()
	service := newCardService(cardRepo, userRepo, langRepo)

	cardRepo.On("FindById", mock.Anything).Return(&model.Card{})

	// act
	card, err := service.GetById(1)

	// assert
	assert.NotNil(t, card)
	assert.Nil(t, err)
}

func Test_Card_ShouldNotGetById(t *testing.T) {
	// arrange
	cardRepo := newMockCardRepository()
	userRepo := newMockUserRepository()
	langRepo := newMockLanguageRepository()
	service := newCardService(cardRepo, userRepo, langRepo)

	cardRepo.On("FindById", mock.Anything).Return(nil)

	// act
	card, err := service.GetById(1)

	// assert
	assert.Nil(t, card)
	assert.NotNil(t, err)
}

func Test_Card_ShouldGetByQuery(t *testing.T) {
	// arrange
	cardRepo := newMockCardRepository()
	userRepo := newMockUserRepository()
	langRepo := newMockLanguageRepository()
	service := newCardService(cardRepo, userRepo, langRepo)

	cardRepo.On("Query", mock.Anything).Return([]*model.Card{}, 0)
	cardRepo.On("Count").Return(0)

	// act
	cards := service.Query(&query.CardQuery{})

	// assert
	assert.NotNil(t, cards)
}

func Test_Card_ShouldUpdate(t *testing.T) {
	// arrange
	cardRepo := newMockCardRepository()
	userRepo := newMockUserRepository()
	langRepo := newMockLanguageRepository()
	service := newCardService(cardRepo, userRepo, langRepo)

	cardRepo.On("FindById", mock.Anything).Return(&model.Card{})
	cardRepo.On("Update", mock.Anything).Return(nil)

	// act
	card, err := service.Update(&dto.PostCard{
		Name:      "card name",
		Text:      "card text",
		Price:     10,
		Type:      "CT1",
		Language:  "ENG",
		Key:       "key1",
		Expansion: "exp1",
	}, 1)

	// assert
	assert.NotNil(t, card)
	assert.Nil(t, err)
}

func Test_Card_ShouldNotUpdate(t *testing.T) {
	// arrange
	cardRepo := newMockCardRepository()
	userRepo := newMockUserRepository()
	langRepo := newMockLanguageRepository()
	service := newCardService(cardRepo, userRepo, langRepo)

	cardRepo.On("FindById", mock.Anything).Return(&model.Card{})
	cardRepo.On("Update", mock.Anything).Return(errors.New(""))

	// act
	card, err := service.Update(&dto.PostCard{
		Name:      "card name",
		Text:      "card text",
		Price:     10,
		Type:      "CT1",
		Language:  "ENG",
		Key:       "key1",
		Expansion: "exp1",
	}, 1)

	// assert
	assert.Nil(t, card)
	assert.NotNil(t, err)
}

func Test_Card_ShouldNotUpdateCardNotFound(t *testing.T) {
	// arrange
	cardRepo := newMockCardRepository()
	userRepo := newMockUserRepository()
	langRepo := newMockLanguageRepository()
	s := newCardService(cardRepo, userRepo, langRepo)

	cardRepo.On("FindById", mock.Anything).Return(nil)

	// act
	card, err := s.Update(&dto.PostCard{
		Name:      "card name",
		Text:      "card text",
		Price:     10,
		Type:      "CT1",
		Language:  "ENG",
		Key:       "key1",
		Expansion: "exp1",
	}, 1)

	// assert
	assert.Nil(t, card)
	assert.Equal(t, service.ErrCardNotFound, err)
}

func Test_Card_ShouldUpdatePrice(t *testing.T) {
	// arrange
	cardRepo := newMockCardRepository()
	userRepo := newMockUserRepository()
	langRepo := newMockLanguageRepository()
	service := newCardService(cardRepo, userRepo, langRepo)

	cardRepo.On("FindById", mock.Anything).Return(&model.Card{})
	cardRepo.On("UpdatePrice", mock.Anything, mock.Anything).Return(&model.Card{}, nil)

	// act
	card, err := service.UpdatePrice(1, &dto.PriceUpdate{
		NewPrice: 10,
	})

	// assert
	assert.NotNil(t, card)
	assert.Nil(t, err)
}

func Test_Card_ShouldNotUpdatePriceInvalidPrice(t *testing.T) {
	// arrange
	cardRepo := newMockCardRepository()
	userRepo := newMockUserRepository()
	langRepo := newMockLanguageRepository()
	service := newCardService(cardRepo, userRepo, langRepo)

	cardRepo.On("FindById", mock.Anything).Return(&model.Card{})
	cardRepo.On("UpdatePrice", mock.Anything, mock.Anything).Return(&model.Card{}, nil)

	// act
	card, err := service.UpdatePrice(1, &dto.PriceUpdate{
		NewPrice: -1,
	})

	// assert
	assert.Nil(t, card)
	assert.NotNil(t, err)
}

func Test_Card_ShouldNotUpdatePriceCardNotFound(t *testing.T) {
	// arrange
	cardRepo := newMockCardRepository()
	userRepo := newMockUserRepository()
	langRepo := newMockLanguageRepository()
	s := newCardService(cardRepo, userRepo, langRepo)

	// cardRepo.On("FindById", mock.Anything).Return(&model.Card{})
	cardRepo.On("UpdatePrice", mock.Anything, mock.Anything).Return(nil, nil)

	// act
	card, err := s.UpdatePrice(1, &dto.PriceUpdate{
		NewPrice: 10,
	})

	// assert
	assert.Nil(t, card)
	assert.Equal(t, service.ErrCardNotFound, err)
}

func Test_Card_ShouldUpdateInStockAmount(t *testing.T) {
	// arrange
	cardRepo := newMockCardRepository()
	userRepo := newMockUserRepository()
	langRepo := newMockLanguageRepository()
	service := newCardService(cardRepo, userRepo, langRepo)

	cardRepo.On("FindById", mock.Anything).Return(&model.Card{})
	cardRepo.On("UpdateInStockAmount", mock.Anything, mock.Anything).Return(&model.Card{}, nil)

	// act
	card, err := service.UpdateInStockAmount(1, &dto.StockedAmountUpdate{
		NewAmount: 10,
	})

	// assert
	assert.NotNil(t, card)
	assert.Nil(t, err)
}

func Test_Card_ShouldNotUpdateInStockAmountCardNotFound(t *testing.T) {
	// arrange
	cardRepo := newMockCardRepository()
	userRepo := newMockUserRepository()
	langRepo := newMockLanguageRepository()
	s := newCardService(cardRepo, userRepo, langRepo)

	// cardRepo.On("FindById", mock.Anything).Return(&model.Card{})
	cardRepo.On("UpdateInStockAmount", mock.Anything, mock.Anything).Return(nil, nil)

	// act
	card, err := s.UpdateInStockAmount(1, &dto.StockedAmountUpdate{
		NewAmount: 10,
	})

	// assert
	assert.Nil(t, card)
	assert.Equal(t, service.ErrCardNotFound, err)
}

func Test_Card_ShouldGetLanguages(t *testing.T) {
	// arrange
	cardRepo := newMockCardRepository()
	userRepo := newMockUserRepository()
	langRepo := newMockLanguageRepository()
	service := newCardService(cardRepo, userRepo, langRepo)

	langRepo.On("All", mock.Anything).Return([]*model.Language{})

	// act
	langs := service.Languages()

	// assert
	assert.Len(t, langs, 0)
}
