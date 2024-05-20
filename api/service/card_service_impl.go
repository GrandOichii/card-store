package service

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"store.api/dto"
	"store.api/model"
	"store.api/repository"
	"store.api/utility"
)

type CardServiceImpl struct {
	CardService
	cardRepo repository.CardRepository
	userRepo repository.UserRepository
	validate *validator.Validate
}

func NewCardServiceImpl(cardRepo repository.CardRepository, userRepo repository.UserRepository, validate *validator.Validate) *CardServiceImpl {
	return &CardServiceImpl{
		cardRepo: cardRepo,
		userRepo: userRepo,
		validate: validate,
	}
}

func (s *CardServiceImpl) GetAll() []*dto.GetCard {
	return utility.MapSlice(s.cardRepo.FindAll(), func(c *model.Card) *dto.GetCard {
		return dto.NewGetCard(c)
	})
}

func (s *CardServiceImpl) Add(c *dto.CreateCard, posterUsername string) (*dto.GetCard, error) {
	err := s.validate.Struct(c)
	if err != nil {
		return nil, err
	}

	// TODO? add more stuff
	card := c.ToCard()

	poster := s.userRepo.FindByUsername(posterUsername)
	if poster == nil {
		return nil, fmt.Errorf("user with username %s doesn't exist", posterUsername)
	}

	card.PosterId = poster.ID
	err = s.cardRepo.Save(card)
	if err != nil {
		return nil, err
	}

	return dto.NewGetCard(card), nil
}
