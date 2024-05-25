package repository

import (
	"gorm.io/gorm"
	"store.api/config"
	"store.api/model"
)

type CardDbRepository struct {
	db     *gorm.DB
	config *config.Configuration
}

func NewCardDbRepository(db *gorm.DB, config *config.Configuration) *CardDbRepository {
	return &CardDbRepository{
		db:     db,
		config: config,
	}
}

func (r *CardDbRepository) Save(card *model.Card) error {
	create := r.db.Create(card)
	err := create.Error
	if err != nil {
		return err
	}

	return nil
}

func (r *CardDbRepository) FindById(id uint) *model.Card {
	var result model.Card
	find := r.db.
		Preload("CardType").
		Preload("Language").
		First(&result, id)
	if find.Error != nil {
		if find.Error == gorm.ErrRecordNotFound {
			return nil
		}
		panic(find.Error)
	}
	return &result
}

func (r *CardDbRepository) Query(page uint, applyQueryF func(*gorm.DB) *gorm.DB) []*model.Card {
	var result []*model.Card
	db := applyQueryF(r.db)
	pageSize := int(r.config.Db.Cards.PageSize)
	offset := (int(page) - 1) * pageSize
	err := db.
		Offset(offset).
		Limit(pageSize).
		Preload("CardType").
		Preload("Language").
		Find(&result).Error
	if err != nil {
		panic(err)
	}
	return result
}
