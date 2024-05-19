package repository

import (
	"gorm.io/gorm"
	"store.api/config"
	"store.api/model"
)

type CardDbRepository struct {
	CardRepository

	db     *gorm.DB
	config *config.Configuration
}

func NewCardDbRepository(db *gorm.DB, config *config.Configuration) *CardDbRepository {
	return &CardDbRepository{
		db:     db,
		config: config,
	}
}

func (r *CardDbRepository) FindAll() []*model.Card {
	var result []*model.Card
	err := r.db.Find(&result).Error
	if err != nil {
		panic(err)
	}
	return result
}

func (r *CardDbRepository) Save(card *model.Card) error {
	err := r.db.Create(card).Error
	if err != nil {
		return err
	}
	return nil
}
