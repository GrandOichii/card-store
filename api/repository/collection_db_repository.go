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
		Preload("Cards").
		Where("owner_id=?", ownerId).
		Find(&result)

	if find.Error != nil {
		panic(find.Error)
	}
	return result
}
