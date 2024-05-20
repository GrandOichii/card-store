package dto

import "store.api/model"

type GetCard struct {
	Name     string `json:"name"`
	Text     string `json:"text"`
	ImageUrl string `json:"imageUrl"`
}

func NewGetCard(c *model.Card) *GetCard {
	return &GetCard{
		Name:     c.Name,
		Text:     c.Text,
		ImageUrl: c.ImageUrl,
	}
}
