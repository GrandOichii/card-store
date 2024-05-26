package repository

import (
	"fmt"

	"gorm.io/gorm"
	"store.api/cache"
	"store.api/config"
	"store.api/model"
)

type CardDbRepository struct {
	db     *gorm.DB
	config *config.Configuration
	cache  cache.CardCache
}

func errCreatedAndFailedToFindCard(id uint) error {
	return fmt.Errorf("created card with id %d, but failed to fetch it", id)
}

func NewCardDbRepository(db *gorm.DB, config *config.Configuration, cache cache.CardCache) *CardDbRepository {
	return &CardDbRepository{
		db:     db,
		config: config,
		cache:  cache,
	}
}

func (r *CardDbRepository) applyPreloads(db *gorm.DB) *gorm.DB {
	return db.
		Preload("CardType").
		Preload("Language")
}

func (r *CardDbRepository) dbFindById(id uint) *model.Card {
	var result model.Card
	find := r.applyPreloads(r.db).
		First(&result, id)

	if find.Error != nil {
		if find.Error == gorm.ErrRecordNotFound {
			return nil
		}
		panic(find.Error)
	}
	return &result
}

func (r *CardDbRepository) Save(card *model.Card) error {
	create := r.db.Create(card)
	err := create.Error
	if err != nil {
		return err
	}
	result := r.dbFindById(card.ID)
	if result == nil {
		panic(errCreatedAndFailedToFindCard(card.ID))
	}
	*card = *result
	r.cache.Remember(result)
	return nil
}

func (r *CardDbRepository) FindById(id uint) *model.Card {
	cached := r.cache.Get(id)
	if cached != nil {
		return cached
	}
	result := r.dbFindById(id)
	if result == nil {
		return nil
	}
	r.cache.Remember(result)
	return result
}

func (r *CardDbRepository) Query(page uint, applyQueryF func(*gorm.DB) *gorm.DB) []*model.Card {
	var result []*model.Card
	db := applyQueryF(r.db)
	pageSize := int(r.config.Db.Cards.PageSize)
	offset := (int(page) - 1) * pageSize
	err := r.applyPreloads(db).
		Offset(offset).
		Limit(pageSize).
		Find(&result).Error
	if err != nil {
		panic(err)
	}
	return result
}

func (r *CardDbRepository) Update(card *model.Card) error {
	update := r.db.Save(card)
	err := update.Error
	if err != nil {
		return err
	}
	result := r.dbFindById(card.ID)
	if result == nil {
		panic(errCreatedAndFailedToFindCard(card.ID))
	}
	*card = *result
	r.cache.Remember(result)
	return nil
}
