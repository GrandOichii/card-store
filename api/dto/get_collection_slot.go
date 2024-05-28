package dto

import "store.api/model"

type GetCollectionSlot struct {
	CardId uint `json:"cardId"`
	Amount uint `json:"amount"`
}

func NewGetCollectionSlot(card *model.CollectionSlot) *GetCollectionSlot {
	return &GetCollectionSlot{
		CardId: card.CardID,
		Amount: card.Amount,
	}
}
