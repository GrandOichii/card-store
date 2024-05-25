package dto

import "store.api/model"

type PostCard struct {
	Name     string  `json:"name" validate:"required"`
	Text     string  `json:"text" validate:"required"`
	ImageUrl string  `json:"imageUrl"`
	Price    float32 `json:"price" validate:"required,gt=0"`
	Type     string  `json:"type" validate:"required"`
	Language string  `json:"language" validate:"required"`
}

func (c PostCard) ToCard() *model.Card {
	return &model.Card{
		Name:       c.Name,
		Text:       c.Text,
		ImageUrl:   c.ImageUrl,
		Price:      c.Price,
		CardTypeID: c.Type,
		LanguageID: c.Language,
	}
}
