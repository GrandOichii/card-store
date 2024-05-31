package dto

import "store.api/model"

type GetCartSlot struct {
	Amount uint `gorm:"not null" json:"amount"`
	CardId uint `gorm:"not null" json:"cardId"`
}

func NewGetCartSlot(slot *model.CartSlot) *GetCartSlot {
	return &GetCartSlot{
		Amount: slot.Amount,
		CardId: slot.CardID,
	}
}
