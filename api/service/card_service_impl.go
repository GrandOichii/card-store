package service

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"store.api/dto"
	"store.api/model"
	"store.api/repository"
	"store.api/utility"
)

type CardServiceImpl struct {
	CardService
	repo     repository.CardRepository
	validate *validator.Validate
}

func NewCardServiceImpl(repo repository.CardRepository, validate *validator.Validate) *CardServiceImpl {
	return &CardServiceImpl{
		repo:     repo,
		validate: validate,
	}
}

func (s *CardServiceImpl) GetAll() []*dto.GetCard {
	return utility.MapSlice(s.repo.FindAll(), func(c *model.Card) *dto.GetCard {
		return dto.NewGetCard(c)
	})
}

func (s *CardServiceImpl) Add(c *dto.CreateCard) (*dto.GetCard, error) {
	// TODO
	return nil, errors.New("not implemented")
}
