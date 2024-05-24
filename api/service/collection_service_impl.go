package service

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"store.api/dto"
	"store.api/model"
	"store.api/repository"
	"store.api/utility"
)

type CollectionServiceImpl struct {
	colRepo  repository.CollectionRepository
	userRepo repository.UserRepository
	validate *validator.Validate
}

func NewCollectionServiceImpl(colRepo repository.CollectionRepository, userRepo repository.UserRepository, validate *validator.Validate) *CollectionServiceImpl {
	return &CollectionServiceImpl{
		colRepo:  colRepo,
		userRepo: userRepo,
		validate: validate,
	}
}

func (ser *CollectionServiceImpl) GetAll(userId uint) []*dto.GetCollection {
	return utility.MapSlice(
		ser.colRepo.FindByOwnerId(userId),
		func(c *model.Collection) *dto.GetCollection { return dto.NewGetCollection(c) },
	)
}

func (ser *CollectionServiceImpl) Create(col *dto.CreateCollection, userId uint) (*dto.GetCollection, error) {
	err := ser.validate.Struct(col)
	if err != nil {
		return nil, err
	}

	user := ser.userRepo.FindById(userId)
	if user == nil {
		return nil, fmt.Errorf("no user with id %d", userId)
	}

	result := col.ToCollection()
	result.OwnerID = userId

	err = ser.colRepo.Save(result)
	if err != nil {
		return nil, err
	}

	return dto.NewGetCollection((result)), nil
}

func (ser *CollectionServiceImpl) AddCard(newCardSlot *dto.CreateCardSlot, colId uint, userId uint) (*dto.GetCollection, error) {
	err := ser.validate.Struct(newCardSlot)
	if err != nil {
		return nil, err
	}

	user := ser.userRepo.FindById(userId)
	if user == nil {
		return nil, fmt.Errorf("no user with id %d", userId)
	}

	collection := ser.colRepo.FindById(colId)
	if collection == nil {
		return nil, err
	}

	// TODO check if cardslot is already present, if so, just update the value
	cardSlot := newCardSlot.ToCardSlot()
	cardSlot.CollectionID = colId
	collection.Cards = append(collection.Cards, *cardSlot)
	err = ser.colRepo.Update(collection)
	if err != nil {
		return nil, err
	}

	updated := ser.colRepo.FindById(collection.ID)

	return dto.NewGetCollection(updated), nil
}
