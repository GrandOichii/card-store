package repository

import (
	"gorm.io/gorm"
	"store.api/config"
	"store.api/model"
)

type CardKeyDbRepository struct {
	db     *gorm.DB
	config *config.Configuration
	// TODO? add cache
}

func NewCardKeyDbRepository(db *gorm.DB, config *config.Configuration) *CardKeyDbRepository {
	return &CardKeyDbRepository{
		db:     db,
		config: config,
	}
}

func (repo *CardKeyDbRepository) All() []*model.CardKey {
	var result []*model.CardKey
	err := repo.db.Find(&result).Error
	if err != nil {
		panic(err)
	}
	return result
}
