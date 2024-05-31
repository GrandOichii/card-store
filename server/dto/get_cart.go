package dto

import (
	"store.api/model"
	"store.api/utility"
)

type GetCart struct {
	Cards []*GetCartSlot `json:"cards"`
}

func NewGetCart(cart *model.Cart) *GetCart {
	return &GetCart{
		Cards: utility.MapSlice(
			cart.Cards,
			func(c model.CartSlot) *GetCartSlot {
				return NewGetCartSlot(&c)
			},
		),
	}
}
