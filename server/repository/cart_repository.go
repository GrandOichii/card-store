package repository

import "store.api/model"

type CartRepository interface {
	Save(cart *model.Cart) error
	FindSingleByUserId(userId uint) *model.Cart
	Update(*model.Cart) error
	UpdateSlot(slot *model.CartSlot) error
	DeleteSlot(slot *model.CartSlot) error
}
