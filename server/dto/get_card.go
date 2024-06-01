package dto

import "store.api/model"

type GetCard struct {
	ID            uint           `json:"id"`
	Name          string         `json:"name"`
	Text          string         `json:"text"`
	ImageUrl      string         `json:"imageUrl"`
	Price         float32        `json:"price"`
	Type          model.CardType `json:"cardType"`
	Language      model.Language `json:"language"`
	Foiling       model.Foiling  `json:"foiling"`
	Key           string         `json:"key"`
	Expansion     string         `json:"expansion"`
	ExpansionName string         `json:"expansionName"`
	InStockAmount uint           `json:"inStockAmount"`
}

func NewGetCard(c *model.Card) *GetCard {
	return &GetCard{
		ID:            c.ID,
		Name:          c.Name,
		Text:          c.Text,
		ImageUrl:      c.ImageUrl,
		Price:         c.Price,
		Type:          c.CardType,
		Language:      c.Language,
		Key:           c.CardKeyID,
		Expansion:     c.Expansion.ShortName,
		ExpansionName: c.Expansion.FullName,
		InStockAmount: c.InStockAmount,
		Foiling:       c.Foiling,
	}
}
