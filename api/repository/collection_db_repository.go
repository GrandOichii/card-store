package repository

import (
	"fmt"

	"gorm.io/gorm"
	"store.api/config"
	"store.api/model"
)

func errCreatedAndFailedToFindCollection(id uint) error {
	return fmt.Errorf("created collection with id %v, but failed to fetch it", id)
}

type CollectionDbRepository struct {
	db     *gorm.DB
	config *config.Configuration
}

func NewCollectionDbRepository(db *gorm.DB, config *config.Configuration) *CollectionDbRepository {
	return &CollectionDbRepository{
		db:     db,
		config: config,
	}
}

func (repo *CollectionDbRepository) dbFindById(id uint) *model.Collection {
	var result model.Collection
	find := repo.db.
		Preload("Cards.Card.CardType").
		First(&result, id)

	if find.Error != nil {
		if find.Error == gorm.ErrRecordNotFound {
			return nil
		}
		panic(find.Error)
	}
	return &result
}

func (repo *CollectionDbRepository) FindByOwnerId(ownerId uint) []*model.Collection {
	var result []*model.Collection
	find := repo.db.
		Where("owner_id=?", ownerId).
		Find(&result)

	if find.Error != nil {
		panic(find.Error)
	}
	return result
}

func (repo *CollectionDbRepository) Save(col *model.Collection) error {
	err := repo.db.Create(col).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *CollectionDbRepository) FindById(id uint) *model.Collection {
	return repo.dbFindById(id)
}

func (repo *CollectionDbRepository) Update(collection *model.Collection) error {
	update := repo.db.Save(collection)
	if update.Error != nil {
		return update.Error
	}
	result := repo.dbFindById(collection.ID)
	if result == nil {
		panic(errCreatedAndFailedToFindCollection(collection.ID))
	}
	return update.Error
}

func (repo *CollectionDbRepository) UpdateCardSlot(cardSlot *model.CardSlot) error {
	update := repo.db.Save(cardSlot)
	return update.Error
}

func (repo *CollectionDbRepository) DeleteCardSlot(cardSlot *model.CardSlot) error {
	delete := repo.db.Delete(cardSlot)
	return delete.Error
}

func (repo *CollectionDbRepository) Delete(id uint) error {
	delete := repo.db.Delete(&model.Collection{}, id)
	return delete.Error
}
