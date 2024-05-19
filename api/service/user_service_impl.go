package service

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"store.api/dto"
	"store.api/repository"
)

type UserServiceImpl struct {
	UserService

	repo     repository.UserRepository
	validate *validator.Validate
}

func NewUserServiceImpl(repo repository.UserRepository, validate *validator.Validate) *UserServiceImpl {
	return &UserServiceImpl{
		repo:     repo,
		validate: validate,
	}
}

func (s *UserServiceImpl) Register(details *dto.RegisterDetails) (*dto.PrivateUserInfo, error) {
	err := s.validate.Struct(details)
	if err != nil {
		return nil, err
	}

	existingUsername, err := s.repo.FindByUsername(details.Username)
	if err != nil {
		return nil, err
	}

	if existingUsername != nil {
		return nil, fmt.Errorf("username %s is taken", details.Username)
	}

	existingEmail, err := s.repo.FindByEmail(details.Email)
	if err != nil {
		return nil, err
	}

	if existingEmail != nil && existingEmail.Verified {
		return nil, fmt.Errorf("already registered account for email %s", details.Email)
	}

	newUser, err := details.ToUser()
	if err != nil {
		return nil, err
	}

	err = s.repo.Save(newUser)
	if err != nil {
		return nil, err
	}

	return dto.NewPrivateUserInfo(newUser), nil
}
