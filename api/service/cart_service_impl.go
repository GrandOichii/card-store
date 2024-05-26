package service

import (
	"fmt"

	"store.api/dto"
	"store.api/repository"
)

type CartServiceImpl struct {
	userRepo repository.UserRepository
	cartRepo repository.CartRepository
}

func NewCartServiceImpl(cartRepo repository.CartRepository, userRepo repository.UserRepository) *CartServiceImpl {
	return &CartServiceImpl{
		cartRepo: cartRepo,
		userRepo: userRepo,
	}
}

func (ser *CartServiceImpl) Get(userId uint) (*dto.GetCart, error) {
	user := ser.userRepo.FindById(userId)
	if user == nil {
		return nil, fmt.Errorf("no user with id %d", userId)
	}

	cart := ser.cartRepo.FindSingleByUserId(userId)
	return dto.NewGetCart(cart), nil
}
