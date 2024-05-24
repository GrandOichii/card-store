package controller_test

import (
	"github.com/stretchr/testify/mock"
	"store.api/dto"
	"store.api/model"
	"store.api/query"
)

type MockUserService struct {
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
	mock.Mock
}

func createMockCardService() *MockCardService {
	return new(MockCardService)
}

func (ser *MockCardService) Add(c *dto.CreateCard, posterId uint) (*dto.GetCard, error) {
	args := ser.Called(c, posterId)
	switch card := args.Get(0).(type) {
	case *dto.GetCard:
		return card, args.Error(1)
	case nil:
		return nil, args.Error(1)
	}
	return nil, args.Error(1)
}

func (ser *MockCardService) GetById(id uint) (*dto.GetCard, error) {
	args := ser.Called(id)
	switch card := args.Get(0).(type) {
	case *dto.GetCard:
		return card, args.Error(1)
	case nil:
		return nil, args.Error(1)
	}
	return nil, args.Error(1)
}

func (ser *MockCardService) Query(query *query.CardQuery) []*dto.GetCard {
	args := ser.Called(query)
	return args.Get(0).([]*dto.GetCard)
}

type MockCollectionService struct {
	mock.Mock
}

func createMockCollectionService() *MockCollectionService {
	return new(MockCollectionService)
}

func (ser *MockCollectionService) GetAll(userId uint) []*dto.GetCollection {
	args := ser.Called(userId)
	return args.Get(0).([]*dto.GetCollection)
}

func (ser *MockCollectionService) Create(c *dto.CreateCollection, userId uint) (*dto.GetCollection, error) {
	args := ser.Called(c, userId)
	switch col := args.Get(0).(type) {
	case *dto.GetCollection:
		return col, args.Error(1)
	case nil:
		return nil, args.Error(1)
	}
	return nil, args.Error(1)
}

func (ser *MockCollectionService) AddCard(cs *dto.CreateCardSlot, colId uint, userId uint) (*dto.GetCollection, error) {
	args := ser.Called(cs, colId, userId)
	switch col := args.Get(0).(type) {
	case *dto.GetCollection:
		return col, args.Error(1)
	case nil:
		return nil, args.Error(1)
	}
	return nil, args.Error(1)
}

func (ser *MockCollectionService) GetById(id uint, userId uint) (*dto.GetCollection, error) {
	args := ser.Called(id, userId)
	switch col := args.Get(0).(type) {
	case *dto.GetCollection:
		return col, args.Error(1)
	case nil:
		return nil, args.Error(1)
	}
	return nil, args.Error(1)
}

// ! duplicated from test/service/mocks_test.go
type MockUserRepository struct {
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

func (m *MockUserRepository) FindById(id uint) *model.User {
	args := m.Called(id)
	switch user := args.Get(0).(type) {
	case *model.User:
		return user
	case nil:
		return nil
	}
	return nil
}
