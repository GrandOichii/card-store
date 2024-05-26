package dto

import "store.api/model"

type GetCollectionSlot struct {
	Card   *GetCard `json:"card"`
	Amount uint     `json:"amount"`
}

func NewGetCollectionSlot(card *model.CollectionSlot) *GetCollectionSlot {
	return &GetCollectionSlot{
		Card:   NewGetCard(&card.Card),
		Amount: card.Amount,
	}
}
