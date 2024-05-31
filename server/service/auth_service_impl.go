package service

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"store.api/dto"
	"store.api/model"
	"store.api/repository"
	"store.api/security"
)

type AuthServiceImpl struct {
	userRepo repository.UserRepository
	cartRepo repository.CartRepository
	validate *validator.Validate
}

func NewAuthServiceImpl(userRepo repository.UserRepository, cartRepo repository.CartRepository, validate *validator.Validate) *AuthServiceImpl {
	return &AuthServiceImpl{
		userRepo: userRepo,
		cartRepo: cartRepo,
		validate: validate,
	}
}

func (s *AuthServiceImpl) Register(details *dto.RegisterDetails) error {
	err := s.validate.Struct(details)
	if err != nil {
		return err
	}

	existingUsername := s.userRepo.FindByUsername(details.Username)

	if existingUsername != nil {
		return fmt.Errorf("username %s is taken", details.Username)
	}

	existingEmail := s.userRepo.FindByEmail(details.Email)

	if existingEmail != nil && existingEmail.Verified {
		return fmt.Errorf("already registered account for email %s", details.Email)
	}

	newUser, err := details.ToUser()
	if err != nil {
		return err
	}

	// save user
	err = s.userRepo.Save(newUser)
	if err != nil {
		return err
	}

	cart := &model.Cart{
		UserID: newUser.ID,
	}
	err = s.cartRepo.Save(cart)
	if err != nil {
		// shouldn't ever happen
		panic(err)
	}

	return nil
}

func (s *AuthServiceImpl) Login(user *dto.LoginDetails) (*dto.PrivateUserInfo, error) {
	if err := s.validate.Struct(user); err != nil {
		return nil, err
	}

	existing := s.userRepo.FindByUsername(user.Username)

	if existing == nil {
		return nil, errors.New("incorrect username or password")
	}

	if !security.CheckPasswordHash(user.Password, existing.PasswordHash) {
		return nil, errors.New("incorrect username or password")
	}

	return dto.NewPrivateUserInfo(existing), nil
}
