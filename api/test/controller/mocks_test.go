package controllers_test

import (
	"github.com/stretchr/testify/mock"
	"store.api/dto"
	"store.api/service"
)

type MockUserService struct {
	service.UserService
	mock.Mock
}

func createMockUserService() *MockUserService {
	return new(MockUserService)
}

func (ser *MockUserService) Login(user *dto.LoginDetails) (*dto.PrivateUserInfo, error) {
	args := ser.Called(user)
	switch user := args.Get(0).(type) {
	case *dto.PrivateUserInfo:
		return user, args.Error(1)
	case nil:
		return nil, args.Error(1)
	}
	return nil, args.Error(1)
}

func (ser *MockUserService) Register(user *dto.RegisterDetails) error {
	args := ser.Called(user)
	return args.Error(0)
}
