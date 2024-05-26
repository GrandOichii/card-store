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

func (r *CartDbRepository) findById(id uint) *model.Cart {
	var result model.Cart
	find := r.db.
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

func (r *CartDbRepository) Update(cart *model.Cart) error {
	update := r.db.Save(cart)
	if update.Error != nil {
		return update.Error
	}
	// result := repo.dbFindById(collection.ID)
	// if result == nil {
	// 	panic(errCreatedAndFailedToFindCollection(collection.ID))
	// }
	// repo.cache.Remember(result)
	return nil
}

func (r *CartDbRepository) UpdateSlot(slot *model.CartSlot) error {
	update := r.db.Save(slot)
	if update.Error != nil {
		return update.Error
	}
	// updated := repo.dbFindById(slot.CollectionID)
	// repo.cache.Remember(updated)
	return nil
}

func (r *CartDbRepository) DeleteSlot(slot *model.CartSlot) error {
	delete := r.db.Delete(slot)
	if delete.Error != nil {
		return delete.Error
	}
	// updated := r.findById(slot.CartID)
	// r.cache.Remember(updated)
	return nil
}
