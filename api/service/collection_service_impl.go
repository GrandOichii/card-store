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

func (ser *CollectionServiceImpl) Create(col *dto.CreateCollection, userId uint) (*dto.GetCollection, error) {
	err := ser.validate.Struct(col)
	if err != nil {
		return nil, err
	}

	// TODO check user id

	result := col.ToCollection()
	result.OwnerID = userId

	err = ser.repo.Save(result)
	if err != nil {
		return nil, err
	}

	return dto.NewGetCollection((result)), nil
}

func (ser *CollectionServiceImpl) AddCard(newCardSlot *dto.CreateCardSlot, colId uint, userId uint) (*dto.GetCollection, error) {
	// TODO check if cardslot is already present, if so, just update the value
	err := ser.validate.Struct(newCardSlot)
	if err != nil {
		return nil, err
	}

	// TODO check user id

	collection := ser.repo.FindById(colId)
	if collection == nil {
		return nil, err
	}

	cardSlot := newCardSlot.ToCardSlot()
	cardSlot.CollectionID = colId
	collection.Cards = append(collection.Cards, *cardSlot)
	err = ser.repo.Update(collection)
	if err != nil {
		return nil, err
	}

	updated := ser.repo.FindById(collection.ID)

	return dto.NewGetCollection(updated), nil
}
