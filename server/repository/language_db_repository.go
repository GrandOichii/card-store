package repository

import (
	"gorm.io/gorm"
	"store.api/config"
	"store.api/model"
)

type LanguageDbRepository struct {
	db     *gorm.DB
	config *config.Configuration
	// TODO add cache
}

func NewLanguageDbRepository(db *gorm.DB, config *config.Configuration) *LanguageDbRepository {
	return &LanguageDbRepository{
		db:     db,
		config: config,
	}
}

func (repo *LanguageDbRepository) All() []*model.Language {
	var result []*model.Language
	err := repo.db.Find(&result).Error
	if err != nil {
		panic(err)
	}
	return result
}
