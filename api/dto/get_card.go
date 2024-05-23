package dto

import "store.api/model"

type GetCard struct {
	ID       uint           `json:"id"`
	Name     string         `json:"name"`
	Text     string         `json:"text"`
	ImageUrl string         `json:"imageUrl"`
	Price    float32        `json:"price"`
	Type     model.CardType `json:"cardType"`
}

func NewGetCard(c *model.Card) *GetCard {
	return &GetCard{
		ID:       c.ID,
		Name:     c.Name,
		Text:     c.Text,
		ImageUrl: c.ImageUrl,
		Price:    c.Price,
		Type:     c.CardType,
	}
}
