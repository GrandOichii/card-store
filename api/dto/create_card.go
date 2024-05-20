package dto

import "store.api/model"

type CreateCard struct {
	Name     string `json:"name" validate:"required"`
	Text     string `json:"text" validate:"required"`
	ImageUrl string `json:"imageUrl"`
}

func (c CreateCard) ToCard() *model.Card {
	return &model.Card{
		Name:     c.Name,
		Text:     c.Text,
		ImageUrl: c.ImageUrl,
	}
}
