package dto

import "store.api/model"

type CreateCardSlot struct {
	CardId uint `json:"cardId" validate:"required"`
	Amount uint `json:"amount" validate:"required,gt=0"`
}

func (c *CreateCardSlot) ToCardSlot() *model.CardSlot {
	return &model.CardSlot{
		Amount: c.Amount,
		CardID: c.CardId,
	}
}
