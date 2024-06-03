package repository

import (
	"gorm.io/gorm"
	"store.api/cache"
	"store.api/config"
	"store.api/model"
)

type ExpansionDbRepository struct {
	db     *gorm.DB
	config *config.Configuration
	cache  cache.ExpansionCache
}

func NewExpansionDbRepository(db *gorm.DB, config *config.Configuration, cache cache.ExpansionCache) *ExpansionDbRepository {
	return &ExpansionDbRepository{
		db:     db,
		config: config,
		cache:  cache,
	}
}

func (repo *ExpansionDbRepository) All() []*model.Expansion {
	cached := repo.cache.Get()
	if cached != nil {
		return cached
	}

	var result []*model.Expansion
	err := repo.db.Find(&result).Error
	if err != nil {
		panic(err)
	}

	repo.cache.Remember(result)
	return result
}
