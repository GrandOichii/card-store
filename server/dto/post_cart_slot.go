package dto

import (
	"fmt"

	"store.api/model"
)

type PostCartSlot struct {
	CardId uint `json:"cardId" validate:"required"`
	Amount int  `json:"amount" validate:"required"`
}

func (s *PostCartSlot) ToCartSlot() (*model.CartSlot, error) {
	if s.Amount <= 0 {
		return nil, fmt.Errorf("%d is not a valid amount number for a cart slot", s.Amount)
	}
	return &model.CartSlot{
		Amount: uint(s.Amount),
		CardID: s.CardId,
	}, nil
}
