package repository

import (
	"gorm.io/gorm"
	"store.api/config"
	"store.api/model"
)

type CartDbRepository struct {
	db     *gorm.DB
	config *config.Configuration
	// cache  cache.CartCache
}

func NewCartDbRepository(db *gorm.DB, config *config.Configuration) *CartDbRepository {
	return &CartDbRepository{
		db:     db,
		config: config,
	}
}

func (r *CartDbRepository) Save(cart *model.Cart) error {
	create := r.db.Create(cart)
	err := create.Error
	if err != nil {
		return err
	}
	return nil
}

func (r *CartDbRepository) FindSingleByUserId(userId uint) *model.Cart {
	var result model.Cart

	find := r.db.
		Preload("Cards").
		Where("user_id=?", userId).
		Find(&result)

	if find.Error != nil {
		// TODO? should this be here or in service
		// shouldn't ever happen
		panic(find.Error)
	}

	return &result
}
