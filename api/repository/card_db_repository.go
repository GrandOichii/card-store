package repository

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"store.api/cache"
	"store.api/config"
	"store.api/model"
	"store.api/query"
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

func (r *CardDbRepository) Query(query *query.CardQuery) []*model.Card {
	var result []*model.Card

	db := r.applyQuery(query, r.db)
	pageSize := int(r.config.Db.Cards.PageSize)
	offset := (int(query.Page) - 1) * pageSize
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

func (repo *CardDbRepository) applyQuery(q *query.CardQuery, d *gorm.DB) *gorm.DB {
	result := d.Where("LOWER(name) like ?", "%"+strings.ToLower(q.Name)+"%")
	if len(q.Type) > 0 {
		result = result.Where("card_type_id=?", q.Type)
	}
	if len(q.Language) > 0 {
		result = result.Where("language_id=?", q.Language)
	}
	if q.MaxPrice != -1 {
		result = result.Where("price < ?", q.MaxPrice)
	}
	if q.MinPrice != -1 {
		result = result.Where("price > ?", q.MinPrice)
	}
	return result
}
