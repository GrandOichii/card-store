package service_test

import (
	"github.com/stretchr/testify/mock"
	"store.api/model"
	"store.api/query"
)

type MockUserRepository struct {
	mock.Mock
}

func newMockUserRepository() *MockUserRepository {
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

func newMockCardRepository() *MockCardRepository {
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

func (m *MockCardRepository) Query(query *query.CardQuery) ([]*model.Card, int64) {
	args := m.Called(query)
	return args.Get(0).([]*model.Card), int64(args.Int(1))
}

func (m *MockCardRepository) Update(c *model.Card) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockCardRepository) UpdatePrice(id uint, newPrice float32) (*model.Card, error) {
	args := m.Called(id, newPrice)
	switch card := args.Get(0).(type) {
	case *model.Card:
		return card, args.Error(1)
	case nil:
		return nil, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockCardRepository) UpdateInStockAmount(id uint, newPrice uint) (*model.Card, error) {
	args := m.Called(id, newPrice)
	switch card := args.Get(0).(type) {
	case *model.Card:
		return card, args.Error(1)
	case nil:
		return nil, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockCardRepository) Count() int64 {
	args := m.Called()
	return int64(args.Int(0))
}

type MockCollectionRepository struct {
	mock.Mock
}

func newMockCollectionRepository() *MockCollectionRepository {
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

func (m *MockCollectionRepository) UpdateSlot(slot *model.CollectionSlot) error {
	args := m.Called(slot)
	return args.Error(0)
}

func (m *MockCollectionRepository) DeleteSlot(slot *model.CollectionSlot) error {
	args := m.Called(slot)
	return args.Error(0)
}

func (m *MockCollectionRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockCartRepository struct {
	mock.Mock
}

func newMockCartRepository() *MockCartRepository {
	return new(MockCartRepository)
}

func (m *MockCartRepository) Save(cart *model.Cart) error {
	args := m.Called(cart)
	return args.Error(0)
}

func (m *MockCartRepository) FindSingleByUserId(userId uint) *model.Cart {
	args := m.Called(userId)
	return args.Get(0).(*model.Cart)
}

func (m *MockCartRepository) Update(c *model.Cart) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockCartRepository) UpdateSlot(slot *model.CartSlot) error {
	args := m.Called(slot)
	return args.Error(0)
}

func (m *MockCartRepository) DeleteSlot(slot *model.CartSlot) error {
	args := m.Called(slot)
	return args.Error(0)
}

type MockLanguageRepository struct {
	mock.Mock
}

func newMockLanguageRepository() *MockLanguageRepository {
	return new(MockLanguageRepository)
}

func (m *MockLanguageRepository) All() []*model.Language {
	args := m.Called()
	return args.Get(0).([]*model.Language)
}

type MockExpansionRepository struct {
	mock.Mock
}

func newMockExpansionRepository() *MockExpansionRepository {
	return new(MockExpansionRepository)
}

func (m *MockExpansionRepository) All() []*model.Expansion {
	args := m.Called()
	return args.Get(0).([]*model.Expansion)
}
