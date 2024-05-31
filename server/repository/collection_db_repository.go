package repository

import (
	"fmt"

	"gorm.io/gorm"
	"store.api/cache"
	"store.api/config"
	"store.api/model"
)

func errCreatedAndFailedToFindCollection(id uint) error {
	return fmt.Errorf("created collection with id %v, but failed to fetch it", id)
}

type CollectionDbRepository struct {
	db     *gorm.DB
	config *config.Configuration
	cache  cache.CollectionCache
}

func NewCollectionDbRepository(db *gorm.DB, config *config.Configuration, cache cache.CollectionCache) *CollectionDbRepository {
	return &CollectionDbRepository{
		db:     db,
		config: config,
		cache:  cache,
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
	existing := repo.cache.Get(id)
	if existing != nil {
		return existing
	}
	result := repo.dbFindById(id)
	if result == nil {
		return nil
	}
	repo.cache.Remember(result)
	return result
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
	repo.cache.Remember(result)
	return nil
}

func (repo *CollectionDbRepository) UpdateSlot(slot *model.CollectionSlot) error {
	update := repo.db.Save(slot)
	if update.Error != nil {
		return update.Error
	}
	updated := repo.dbFindById(slot.CollectionID)
	repo.cache.Remember(updated)
	return nil
}

func (repo *CollectionDbRepository) DeleteSlot(slot *model.CollectionSlot) error {
	delete := repo.db.Delete(slot)
	if delete.Error != nil {
		return delete.Error
	}
	updated := repo.dbFindById(slot.CollectionID)
	repo.cache.Remember(updated)
	return nil
}

func (repo *CollectionDbRepository) Delete(id uint) error {
	delete := repo.db.Delete(&model.Collection{}, id)
	if delete.Error != nil {
		return delete.Error
	}
	repo.cache.Forget(id)
	return nil
}
