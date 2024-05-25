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

	if !user.Verified {
		return nil, ErrNotVerified
	}

	result := col.ToCollection()
	result.OwnerID = userId

	err = ser.colRepo.Save(result)
	if err != nil {
		return nil, err
	}

	return dto.NewGetCollection((result)), nil
}

func (ser *CollectionServiceImpl) EditCard(newCardSlot *dto.PostCardSlot, colId uint, userId uint) (*dto.GetCollection, error) {
	err := ser.validate.Struct(newCardSlot)
	if err != nil {
		return nil, err
	}

	user := ser.userRepo.FindById(userId)
	if user == nil {
		return nil, fmt.Errorf("no user with id %d", userId)
	}

	if !user.Verified {
		return nil, ErrNotVerified
	}

	collection, err := ser.getById(colId, userId)
	if err != nil {
		return nil, err
	}

	added := false
	for _, slot := range collection.Cards {
		if slot.CardID == newCardSlot.CardId {
			added = true
			slot.Amount += uint(newCardSlot.Amount)
			if slot.Amount <= 0 {
				err = ser.colRepo.DeleteCardSlot(&slot)
				if err != nil {
					return nil, err
				}
				break
			}
			err = ser.colRepo.UpdateCardSlot(&slot)
			if err != nil {
				return nil, err
			}
			break
		}
	}

	if !added {
		cardSlot, err := newCardSlot.ToCardSlot()
		if err != nil {
			return nil, err
		}
		cardSlot.CollectionID = colId
		collection.Cards = append(collection.Cards, *cardSlot)
		err = ser.colRepo.Update(collection)
		if err != nil {
			return nil, err
		}
	}

	updated := ser.colRepo.FindById(collection.ID)

	return dto.NewGetCollection(updated), nil
}

func (ser *CollectionServiceImpl) GetById(id uint, userId uint) (*dto.GetCollection, error) {
	result, err := ser.getById(id, userId)
	if err != nil {
		return nil, err
	}
	return dto.NewGetCollection(result), nil
}

func (ser *CollectionServiceImpl) Delete(id uint, userId uint) error {
	result, err := ser.getById(id, userId)
	if err != nil {
		return err
	}
	return ser.colRepo.Delete(result.ID)
}

func (ser *CollectionServiceImpl) getById(id uint, userId uint) (*model.Collection, error) {
	result := ser.colRepo.FindById(id)
	if result == nil || result.OwnerID != userId {
		return nil, fmt.Errorf("no collection with id %v", id)
	}
	return result, nil
}
