package service

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"store.api/cache"
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
	cache    cache.CardCache
}

func NewCardServiceImpl(cardRepo repository.CardRepository, userRepo repository.UserRepository, validate *validator.Validate, cache cache.CardCache) *CardServiceImpl {
	return &CardServiceImpl{
		cardRepo: cardRepo,
		userRepo: userRepo,
		validate: validate,
		cache:    cache,
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

	result := dto.NewGetCard(card)

	err = s.cache.Remember(result)
	if err != nil {
		panic(err)
	}

	return result, nil
}

func (s *CardServiceImpl) GetById(id uint) (*dto.GetCard, error) {
	remembered, err := s.cache.Get(id)
	if err != nil {
		panic(err)
	}
	if remembered != nil {
		return remembered, nil
	}

	card := s.cardRepo.FindById(id)
	if card == nil {
		return nil, fmt.Errorf("no card with id %d", id)
	}
	result := dto.NewGetCard(card)

	err = s.cache.Remember(result)
	if err != nil {
		panic(err)
	}

	return result, nil
}

func (s *CardServiceImpl) Query(query *query.CardQuery) []*dto.GetCard {
	// TODO? add cache?
	// TODO move to a more text-search specific service
	applyQueryF := query.ApplyQueryF()
	cards := s.cardRepo.Query(query.Page, applyQueryF)

	return utility.MapSlice(cards, func(c *model.Card) *dto.GetCard {
		return dto.NewGetCard(c)
	})
}
