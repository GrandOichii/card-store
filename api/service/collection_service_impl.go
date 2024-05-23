package service

import (
	"github.com/go-playground/validator/v10"
	"store.api/dto"
	"store.api/model"
	"store.api/repository"
	"store.api/utility"
)

type CollectionServiceImpl struct {
	repo     repository.CollectionRepository
	validate *validator.Validate
}

func NewCollectionServiceImpl(repo repository.CollectionRepository, validate *validator.Validate) *CollectionServiceImpl {
	return &CollectionServiceImpl{
		repo:     repo,
		validate: validate,
	}
}

func (ser *CollectionServiceImpl) GetAll(userId uint) []*dto.GetCollection {
	return utility.MapSlice(
		ser.repo.FindByOwnerId(userId),
		func(c *model.Collection) *dto.GetCollection { return dto.NewGetCollection(c) },
	)
}
