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

func (r *CardDbRepository) FindById(id uint) *model.Card {
	var result model.Card
	find := r.db.First(&result, id)
	if find.Error != nil {
		if find.Error == gorm.ErrRecordNotFound {
			return nil
		}
		panic(find.Error)
	}
	return &result
}

func (r *CardDbRepository) FindByType(cType string) ([]*model.Card, error) {
	var result []*model.Card
	find := r.db.Where("card_type_id=?", cType).Find(&result)
	// TODO dont think it will ever throw an error
	if find.Error != nil {
		return nil, find.Error
	}
	return result, nil
}
