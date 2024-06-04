package service

import (
	"store.api/dto"
	"store.api/repository"
)

// TODO add tests

type UserServiceImpl struct {
	userRepo repository.UserRepository
}

func NewUserServiceImpl(userRepo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

func (ser *UserServiceImpl) ById(id uint) (*dto.PrivateUserInfo, error) {
	user := ser.userRepo.FindById(id)
	if user == nil {
		return nil, ErrUserNotFound
	}

	return dto.NewPrivateUserInfo(user), nil
}
