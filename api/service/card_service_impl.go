package service

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"store.api/dto"
	"store.api/model"
	"store.api/query"
	"store.api/repository"
	"store.api/utility"
)

type CardServiceImpl struct {
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

func (s *CardServiceImpl) Add(c *dto.CreateCard, posterId uint) (*dto.GetCard, error) {
	err := s.validate.Struct(c)
	if err != nil {
		return nil, err
	}

	// TODO add more stuff
	card := c.ToCard()

	poster := s.userRepo.FindById(posterId)
	if poster == nil {
		return nil, fmt.Errorf("user with id %d doesn't exist", posterId)
	}

	card.PosterID = poster.ID
	err = s.cardRepo.Save(card)
	if err != nil {
		return nil, err
	}

	return dto.NewGetCard(card), nil
}

func (s *CardServiceImpl) GetById(id uint) (*dto.GetCard, error) {
	result := s.cardRepo.FindById(id)
	if result == nil {
		return nil, fmt.Errorf("no card with id %d", id)
	}
	return dto.NewGetCard(result), nil
}

func (s *CardServiceImpl) Query(query *query.CardQuery) []*dto.GetCard {
	applyQueryF := query.ApplyQueryF()
	cards := s.cardRepo.Query(applyQueryF)

	return utility.MapSlice(cards, func(c *model.Card) *dto.GetCard {
		return dto.NewGetCard(c)
	})
}
