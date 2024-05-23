package dto

import "store.api/model"

type GetCardSlot struct {
	Card   *GetCard `json:"card"`
	Amount uint     `json:"amount"`
}

func NewGetCardSlot(card *model.CardSlot) *GetCardSlot {
	return &GetCardSlot{
		Card:   NewGetCard(&card.Card),
		Amount: card.Amount,
	}
}
