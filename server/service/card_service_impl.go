package service

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"store.api/config"
	"store.api/dto"
	"store.api/model"
	"store.api/query"
	"store.api/repository"
	"store.api/utility"
)

type CardServiceImpl struct {
	config *config.Configuration

	cardRepo      repository.CardRepository
	userRepo      repository.UserRepository
	langRepo      repository.LanguageRepository
	expansionRepo repository.ExpansionRepository
	validate      *validator.Validate
}

func NewCardServiceImpl(config *config.Configuration, cardRepo repository.CardRepository, userRepo repository.UserRepository, langRepo repository.LanguageRepository, expansionRepo repository.ExpansionRepository, validate *validator.Validate) *CardServiceImpl {
	return &CardServiceImpl{
		config: config,

		cardRepo:      cardRepo,
		userRepo:      userRepo,
		langRepo:      langRepo,
		expansionRepo: expansionRepo,
		validate:      validate,
	}
}

func (s *CardServiceImpl) Add(c *dto.PostCard, posterId uint) (*dto.GetCard, error) {
	err := s.validate.Struct(c)
	if err != nil {
		return nil, err
	}

	card := c.ToCard()

	poster := s.userRepo.FindById(posterId)
	if poster == nil {
		return nil, ErrUserNotFound
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

func (s *CardServiceImpl) Query(query *query.CardQuery) *CardQueryResult {
	// TODO move to a more text-search specific service
	cards, count := s.cardRepo.Query(query)

	mapped := utility.MapSlice(cards, func(c *model.Card) *dto.GetCard {
		return dto.NewGetCard(c)
	})

	return &CardQueryResult{
		Cards:      mapped,
		TotalCount: count,
		PerPage:    s.config.Db.Cards.PageSize,
	}
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

func (s *CardServiceImpl) UpdatePrice(id uint, update *dto.PriceUpdate) (*dto.GetCard, error) {
	if update.NewPrice <= 0 {
		return nil, fmt.Errorf("card price can't be %f", update.NewPrice)
	}
	result, err := s.cardRepo.UpdatePrice(id, update.NewPrice)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, ErrCardNotFound
	}
	return dto.NewGetCard(result), nil
}

func (s *CardServiceImpl) UpdateInStockAmount(id uint, update *dto.StockedAmountUpdate) (*dto.GetCard, error) {
	result, err := s.cardRepo.UpdateInStockAmount(id, update.NewAmount)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, ErrCardNotFound
	}
	return dto.NewGetCard(result), nil
}

func (s *CardServiceImpl) Languages() []*model.Language {
	result := s.langRepo.All()
	return result
}

func (s *CardServiceImpl) Expansions() []*model.Expansion {
	result := s.expansionRepo.All()
	return result
}
