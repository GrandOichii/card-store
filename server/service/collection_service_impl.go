package service

import (
	"github.com/go-playground/validator/v10"
	"store.api/dto"
	"store.api/model"
	"store.api/repository"
	"store.api/utility"
)

type CollectionServiceImpl struct {
	colRepo  repository.CollectionRepository
	userRepo repository.UserRepository
	cardRepo repository.CardRepository
	validate *validator.Validate
}

func NewCollectionServiceImpl(colRepo repository.CollectionRepository, userRepo repository.UserRepository, cardRepo repository.CardRepository, validate *validator.Validate) *CollectionServiceImpl {
	return &CollectionServiceImpl{
		colRepo:  colRepo,
		userRepo: userRepo,
		cardRepo: cardRepo,
		validate: validate,
	}
}

func (ser *CollectionServiceImpl) GetAll(userId uint) []*dto.GetCollection {
	return utility.MapSlice(
		ser.colRepo.FindByOwnerId(userId),
		func(c *model.Collection) *dto.GetCollection { return dto.NewGetCollection(c) },
	)
}

func (ser *CollectionServiceImpl) Create(col *dto.PostCollection, userId uint) (*dto.GetCollection, error) {
	err := ser.validate.Struct(col)
	if err != nil {
		return nil, err
	}

	user := ser.userRepo.FindById(userId)
	if user == nil {
		return nil, ErrUserNotFound
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

func (ser *CollectionServiceImpl) EditSlot(newCollectionSlot *dto.PostCollectionSlot, colId uint, userId uint) (*dto.GetCollection, error) {
	err := ser.validate.Struct(newCollectionSlot)
	if err != nil {
		return nil, err
	}

	user := ser.userRepo.FindById(userId)
	if user == nil {
		return nil, ErrUserNotFound
	}

	if !user.Verified {
		return nil, ErrNotVerified
	}

	collection, err := ser.getById(colId, userId)
	if err != nil {
		return nil, err
	}

	card := ser.cardRepo.FindById(newCollectionSlot.CardId)
	if card == nil {
		return nil, ErrCardNotFound
	}

	added := false
	for _, slot := range collection.Cards {
		if slot.CardID == newCollectionSlot.CardId {
			added = true
			slot.Amount += uint(newCollectionSlot.Amount)
			if slot.Amount <= 0 {
				err = ser.colRepo.DeleteSlot(&slot)
				if err != nil {
					return nil, err
				}
				break
			}
			err = ser.colRepo.UpdateSlot(&slot)
			if err != nil {
				return nil, err
			}
			break
		}
	}

	if !added {
		collectionSlot, err := newCollectionSlot.ToCollectionSlot()
		if err != nil {
			return nil, err
		}
		collectionSlot.CollectionID = colId
		collection.Cards = append(collection.Cards, *collectionSlot)
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

func (ser *CollectionServiceImpl) UpdateInfo(newData *dto.PostCollection, id uint, userId uint) (*dto.GetCollection, error) {
	// check validity
	err := ser.validate.Struct(newData)
	if err != nil {
		return nil, err
	}

	// fetch user
	user := ser.userRepo.FindById(id)
	if user == nil {
		return nil, ErrUserNotFound
	}

	if !user.Verified {
		return nil, ErrNotVerified
	}

	existing, err := ser.getById(id, userId)
	if err != nil {
		return nil, err
	}

	// modify collection
	newCollection := newData.ToCollection()
	newCollection.ID = existing.ID
	newCollection.OwnerID = existing.OwnerID
	newCollection.Cards = existing.Cards

	err = ser.colRepo.Update(newCollection)
	if err != nil {
		return nil, err
	}

	return dto.NewGetCollection(newCollection), nil
}

func (ser *CollectionServiceImpl) getById(id uint, userId uint) (*model.Collection, error) {
	result := ser.colRepo.FindById(id)
	if result == nil || result.OwnerID != userId {
		return nil, ErrCollectionNotFound
	}
	return result, nil
}
