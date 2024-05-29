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
	db         *gorm.DB
	config     *config.Configuration
	cardCache  cache.CardCache
	queryCache cache.CardQueryCache
}

func errCreatedAndFailedToFindCard(id uint) error {
	return fmt.Errorf("created card with id %d, but failed to fetch it", id)
}

func NewCardDbRepository(db *gorm.DB, config *config.Configuration, cardCache cache.CardCache, queryCache cache.CardQueryCache) *CardDbRepository {
	return &CardDbRepository{
		db:         db,
		config:     config,
		cardCache:  cardCache,
		queryCache: queryCache,
	}
}

func (r *CardDbRepository) applyPreloads(db *gorm.DB) *gorm.DB {
	return db.
		Preload("CardType").
		Preload("Expansion").
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
	r.cardCache.Remember(result)

	// TODO not tested
	r.queryCache.ForgetAll()
	return nil
}

func (r *CardDbRepository) FindById(id uint) *model.Card {
	cached := r.cardCache.Get(id)
	if cached != nil {
		return cached
	}
	result := r.dbFindById(id)
	if result == nil {
		return nil
	}
	r.cardCache.Remember(result)
	return result
}

func (r *CardDbRepository) Query(query *query.CardQuery) []*model.Card {
	cached := r.queryCache.Get(query.Raw)
	if cached != nil {
		return cached
	}
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

	r.queryCache.Remember(query.Raw, result)

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
	r.cardCache.Remember(result)

	// TODO not tested
	r.queryCache.ForgetAll()

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
	if len(q.Key) > 0 {
		result = result.Where("card_key_id=?", q.Key)
	}
	if len(q.Expansion) > 0 {
		result = result.Where("expansion_id=?", q.Expansion)
	}
	if len(q.Keywords) > 0 {
		// oh boy

		// keywords can contain:
		// v parts of card name: could be one word, could be words not in order, could be parts of words
		// v card type: lowercase card types, like MTG or ygo, also by short name, like magic or yugioh
		// v card language: language symbol or full names: rus, eng, english
		// - tags: special tags that are attached to cards to make searching easier TODO

		// keywords CAN'T contain (for now):
		// - card types: don't see a reason for this
		// - card cost/power/toughness/life/etc: also don't see a reason for this, unless someone wants to build 6cmc tribal
		// - date of printing: why

		// keywords under consideration:
		// - collectors number: could be pretty useful for collectors, but still very niche
		// - author: also for collection purposes

		result = result.Joins("JOIN languages ON cards.language_id = languages.id")
		result = result.Joins("JOIN card_types ON cards.card_type_id = card_types.id")
		result = result.Joins("JOIN card_keys ON cards.card_key_id = card_keys.id")
		result = result.Joins("JOIN expansions ON cards.expansion_id = expansions.id")

		words := strings.Split(q.Keywords, " ")
		fmt.Printf("len(words): %v\n", len(words))
		for _, word := range words {
			w := strings.ToLower(word)

			// ! must be ordered correctly!

			// ? sort words from shortest to longest?

			result = result.
				// language
				Where("(LOWER(language_id) = ?", w).
				Or("LOWER(languages.long_name) = ?", w).
				// type
				Or("LOWER(cards.card_type_id) = ?", w).
				Or("LOWER(card_types.short_name) = ?", w).
				// expansion
				Or("LOWER(expansions.short_name) = ?", w).
				// name
				Or("LOWER(name) like ?", "%"+w+"%").
				// key name
				Or("LOWER(card_keys.eng_name) like ?)", "%"+w+"%")
		}
	}
	return result
}
