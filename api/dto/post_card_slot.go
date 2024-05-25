package dto

import (
	"fmt"

	"store.api/model"
)

type PostCardSlot struct {
	CardId uint `json:"cardId" validate:"required"`
	Amount int  `json:"amount" validate:"required"`
}

func (c *PostCardSlot) ToCardSlot() (*model.CardSlot, error) {
	if c.Amount <= 0 {
		return nil, fmt.Errorf("%d is not a valid amount number for a card slot", c.Amount)
	}
	return &model.CardSlot{
		Amount: uint(c.Amount),
		CardID: c.CardId,
	}, nil
}
