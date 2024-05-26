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

func (s *CardServiceImpl) Add(c *dto.PostCard, posterId uint) (*dto.GetCard, error) {
	err := s.validate.Struct(c)
	if err != nil {
		return nil, err
	}

	// TODO add more stuff
	card := c.ToCard()

	poster := s.userRepo.FindById(posterId)
	if poster == nil {
		return nil, fmt.Errorf("no user with id %d", posterId)
	}

	card.PosterID = poster.ID
	err = s.cardRepo.Save(card)
	if err != nil {
		return nil, err
	}

	return dto.NewGetCard(card), nil
}

func (s *CardServiceImpl) GetById(id uint) (*dto.GetCard, error) {
	card := s.cardRepo.FindById(id)
	if card == nil {
		return nil, ErrCardNotFound
	}
	result := dto.NewGetCard(card)

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

func (s *CardServiceImpl) Update(c *dto.PostCard, cardId uint) (*dto.GetCard, error) {
	err := s.validate.Struct(c)
	if err != nil {
		return nil, err
	}

	newCard := c.ToCard()
	existing := s.cardRepo.FindById(cardId)
	if existing == nil {
		return nil, ErrCardNotFound
	}

	newCard.ID = existing.ID
	newCard.PosterID = existing.PosterID
	err = s.cardRepo.Update(newCard)
	if err != nil {
		return nil, err
	}

	return dto.NewGetCard(newCard), nil
}
