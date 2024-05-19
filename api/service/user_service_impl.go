package service

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"store.api/dto"
	"store.api/repository"
	"store.api/security"
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

func (s *UserServiceImpl) Register(details *dto.RegisterDetails) error {
	err := s.validate.Struct(details)
	if err != nil {
		return err
	}

	existingUsername, err := s.repo.FindByUsername(details.Username)
	if err != nil {
		return err
	}

	if existingUsername != nil {
		return fmt.Errorf("username %s is taken", details.Username)
	}

	existingEmail, err := s.repo.FindByEmail(details.Email)
	if err != nil {
		return err
	}

	if existingEmail != nil && existingEmail.Verified {
		return fmt.Errorf("already registered account for email %s", details.Email)
	}

	newUser, err := details.ToUser()
	if err != nil {
		return err
	}

	err = s.repo.Save(newUser)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserServiceImpl) Login(user *dto.LoginDetails) (*dto.PrivateUserInfo, error) {
	if err := s.validate.Struct(user); err != nil {
		return nil, err
	}

	existing, err := s.repo.FindByUsername(user.Username)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, errors.New("incorrect username or password")
	}

	if !security.CheckPasswordHash(user.Password, existing.PasswordHash) {
		return nil, errors.New("incorrect username or password")
	}

	return dto.NewPrivateUserInfo(existing), nil
}
