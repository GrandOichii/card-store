package controller_test

import (
	"github.com/stretchr/testify/mock"
	"store.api/dto"
	"store.api/model"
	"store.api/repository"
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

type MockCardService struct {
	service.CardService
	mock.Mock
}

func createMockCardService() *MockCardService {
	return new(MockCardService)
}

func (ser *MockCardService) GetAll() []*dto.GetCard {
	args := ser.Called()
	return args.Get(0).([]*dto.GetCard)
}

func (ser *MockCardService) Add(c *dto.CreateCard, posterId string) (*dto.GetCard, error) {
	args := ser.Called(c, posterId)
	switch card := args.Get(0).(type) {
	case *dto.GetCard:
		return card, args.Error(1)
	case nil:
		return nil, args.Error(1)
	}
	return nil, args.Error(1)
}

// ! duplicated from test/service/mocks_test.go
type MockUserRepository struct {
	repository.UserRepository
	mock.Mock
}

func createMockUserRepository() *MockUserRepository {
	return new(MockUserRepository)
}

func (m *MockUserRepository) FindByUsername(username string) *model.User {
	args := m.Called(username)
	switch user := args.Get(0).(type) {
	case *model.User:
		return user
	case nil:
		return nil
	}
	return nil
}

func (m *MockUserRepository) FindByEmail(email string) *model.User {
	args := m.Called(email)
	switch user := args.Get(0).(type) {
	case *model.User:
		return user
	case nil:
		return nil
	}
	return nil
}

func (m *MockUserRepository) Save(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}
