package service

import "store.api/dto"

type CartService interface {
	Get(userId uint) (*dto.GetCart, error)
}
