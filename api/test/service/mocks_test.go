package service_test

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
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

func (m *MockCardRepository) Query(page uint, applyQueryF func(*gorm.DB) *gorm.DB) []*model.Card {
	args := m.Called(page, applyQueryF)
	return args.Get(0).([]*model.Card)
}

type MockCollectionRepository struct {
	mock.Mock
}

func createMockCollectionRepository() *MockCollectionRepository {
	return new(MockCollectionRepository)
}

func (m *MockCollectionRepository) FindByOwnerId(ownerId uint) []*model.Collection {
	args := m.Called(ownerId)
	return args.Get(0).([]*model.Collection)
}

func (m *MockCollectionRepository) Save(c *model.Collection) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockCollectionRepository) FindById(id uint) *model.Collection {
	args := m.Called(id)
	switch col := args.Get(0).(type) {
	case *model.Collection:
		return col
	case nil:
		return nil
	}
	return nil
}

func (m *MockCollectionRepository) Update(c *model.Collection) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockCollectionRepository) UpdateCardSlot(cardSlot *model.CardSlot) error {
	args := m.Called(cardSlot)
	return args.Error(0)
}

func (m *MockCollectionRepository) DeleteCardSlot(cardSlot *model.CardSlot) error {
	args := m.Called(cardSlot)
	return args.Error(0)
}
