package repository

import (
	"gorm.io/gorm"
	"store.api/cache"
	"store.api/config"
	"store.api/model"
)

type LanguageDbRepository struct {
	db     *gorm.DB
	config *config.Configuration
	cache  cache.LanguageCache
}

func NewLanguageDbRepository(db *gorm.DB, config *config.Configuration, cache cache.LanguageCache) *LanguageDbRepository {
	return &LanguageDbRepository{
		db:     db,
		config: config,
		cache:  cache,
	}
}

func (repo *LanguageDbRepository) All() []*model.Language {
	cached := repo.cache.Get()
	if cached != nil {
		return cached
	}

	var result []*model.Language
	err := repo.db.Find(&result).Error
	if err != nil {
		panic(err)
	}

	repo.cache.Remember(result)
	return result
}
