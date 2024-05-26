package repository

import "store.api/model"

type CartRepository interface {
	Save(cart *model.Cart) error
	FindSingleByUserId(userId uint) *model.Cart
}
