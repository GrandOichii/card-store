package dto

import "store.api/model"

type GetCard struct {
	// TODO
	Name string `json:"name"`
	Text string `json:"text"`
}

func NewGetCard(c *model.Card) *GetCard {
	return &GetCard{
		Name: c.Name,
		Text: c.Text,
		// TODO
	}
}
