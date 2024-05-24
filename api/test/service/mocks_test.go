package service_test

import (
	"github.com/stretchr/testify/mock"
	"store.api/model"
)

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

type MockCardRepository struct {
	mock.Mock
}

func createMockCardRepository() *MockCardRepository {
	return new(MockCardRepository)
}

func (m *MockCardRepository) Save(c *model.Card) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockCardRepository) FindById(id uint) *model.Card {
	args := m.Called(id)
	switch card := args.Get(0).(type) {
	case *model.Card:
		return card
	case nil:
		return nil
	}
	return nil
}

func (m *MockCardRepository) FindByType(cType string) []*model.Card {
	args := m.Called(cType)
	return args.Get(0).([]*model.Card)
}
