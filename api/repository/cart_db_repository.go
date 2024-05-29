package repository

import (
	"gorm.io/gorm"
	"store.api/cache"
	"store.api/config"
	"store.api/model"
)

type CartDbRepository struct {
	db     *gorm.DB
	config *config.Configuration
	cache  cache.CartCache
}

func NewCartDbRepository(db *gorm.DB, config *config.Configuration, cache cache.CartCache) *CartDbRepository {
	return &CartDbRepository{
		db:     db,
		config: config,
		cache:  cache,
	}
}

func (r *CartDbRepository) dbFindSingleByUserId(userId uint) *model.Cart {
	var result model.Cart

	find := r.db.
		Preload("Cards").
		Where("user_id=?", userId).
		Find(&result)

	if find.Error != nil {
		// shouldn't ever happen
		panic(find.Error)
	}

	return &result
}

func (r *CartDbRepository) dbFindById(id uint) *model.Cart {
	var result model.Cart
	find := r.db.
		Preload("Cards").
		First(&result, id)

	if find.Error != nil {
		if find.Error == gorm.ErrRecordNotFound {
			return nil
		}
		panic(find.Error)
	}
	return &result
}
func (r *CartDbRepository) Save(cart *model.Cart) error {
	create := r.db.Create(cart)
	err := create.Error
	if err != nil {
		return err
	}
	r.cache.Remember(cart)
	return nil
}

func (r *CartDbRepository) FindSingleByUserId(userId uint) *model.Cart {
	cached := r.cache.Get(userId)
	if cached != nil {
		return cached
	}
	result := r.dbFindSingleByUserId(userId)
	r.cache.Remember(result)
	return result
}

func (r *CartDbRepository) Update(cart *model.Cart) error {
	update := r.db.Save(cart)
	if update.Error != nil {
		return update.Error
	}
	result := r.dbFindSingleByUserId(cart.UserID)
	r.cache.Remember(result)
	return nil
}

func (r *CartDbRepository) UpdateSlot(slot *model.CartSlot) error {
	update := r.db.Save(slot)
	if update.Error != nil {
		return update.Error
	}
	updated := r.dbFindById(slot.CartID)
	r.cache.Remember(updated)
	return nil
}

func (r *CartDbRepository) DeleteSlot(slot *model.CartSlot) error {
	delete := r.db.Delete(slot)
	if delete.Error != nil {
		return delete.Error
	}
	updated := r.dbFindById(slot.CartID)
	r.cache.Remember(updated)
	return nil
}
