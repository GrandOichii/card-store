package dto

import (
	"fmt"

	"store.api/model"
)

type PostCollectionSlot struct {
	CardId uint `json:"cardId" validate:"required"`
	Amount int  `json:"amount" validate:"required"`
}

func (c *PostCollectionSlot) ToCollectionSlot() (*model.CollectionSlot, error) {
	if c.Amount <= 0 {
		return nil, fmt.Errorf("%d is not a valid amount number for a collection slot", c.Amount)
	}
	return &model.CollectionSlot{
		Amount: uint(c.Amount),
		CardID: c.CardId,
	}, nil
}
