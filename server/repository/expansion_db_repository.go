package repository

import (
	"gorm.io/gorm"
	"store.api/config"
	"store.api/model"
)

type ExpansionDbRepository struct {
	db     *gorm.DB
	config *config.Configuration
	// TODO add cache
}

func NewExpansionDbRepository(db *gorm.DB, config *config.Configuration) *ExpansionDbRepository {
	return &ExpansionDbRepository{
		db:     db,
		config: config,
	}
}

func (repo *ExpansionDbRepository) All() []*model.Expansion {
	var result []*model.Expansion
	err := repo.db.Find(&result).Error
	if err != nil {
		panic(err)
	}
	return result
}
