package repository

import (
	"gorm.io/gorm"
	"store.api/config"
	"store.api/model"
)

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

func (repo *CollectionDbRepository) FindByOwnerId(ownerId uint) []*model.Collection {
	var result []*model.Collection
	find := repo.db.
		Preload("Cards.Card.CardType").
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

func (repo *CollectionDbRepository) Update(collection *model.Collection) error {
	update := repo.db.Save(*collection)
	return update.Error
}
