package service

import (
	"errors"

	"store.api/dto"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type CartService interface {
	Get(userId uint) (*dto.GetCart, error)
	EditSlot(userId uint, cartSlot *dto.PostCartSlot) (*dto.GetCart, error)
}
