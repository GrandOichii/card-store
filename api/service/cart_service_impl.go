package service

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"store.api/dto"
	"store.api/repository"
)

type CartServiceImpl struct {
	userRepo repository.UserRepository
	cartRepo repository.CartRepository
	cardRepo repository.CardRepository
	validate *validator.Validate
}

func NewCartServiceImpl(cartRepo repository.CartRepository, userRepo repository.UserRepository, cardRepo repository.CardRepository, validate *validator.Validate) *CartServiceImpl {
	return &CartServiceImpl{
		cartRepo: cartRepo,
		userRepo: userRepo,
		cardRepo: cardRepo,
		validate: validate,
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

func (ser *CartServiceImpl) EditSlot(userId uint, newCartSlot *dto.PostCartSlot) (*dto.GetCart, error) {
	err := ser.validate.Struct(newCartSlot)
	if err != nil {
		return nil, err
	}

	user := ser.userRepo.FindById(userId)
	if user == nil {
		return nil, fmt.Errorf("no user with id %d", userId)
	}

	card := ser.cardRepo.FindById(newCartSlot.CardId)
	if card == nil {
		return nil, ErrCardNotFound
	}

	cart := ser.cartRepo.FindSingleByUserId(userId)

	added := false
	for _, slot := range cart.Cards {
		if slot.CardID == newCartSlot.CardId {
			added = true
			slot.Amount += uint(newCartSlot.Amount)
			if slot.Amount <= 0 {
				err = ser.cartRepo.DeleteSlot(&slot)
				if err != nil {
					return nil, err
				}
				break
			}
			err = ser.cartRepo.UpdateSlot(&slot)
			if err != nil {
				return nil, err
			}
			break
		}
	}

	if !added {
		collectionSlot, err := newCartSlot.ToCartSlot()
		if err != nil {
			return nil, err
		}
		cart.Cards = append(cart.Cards, *collectionSlot)
		err = ser.cartRepo.Update(cart)
		if err != nil {
			return nil, err
		}
	}

	updated := ser.cartRepo.FindSingleByUserId(userId)
	return dto.NewGetCart(updated), nil
}
