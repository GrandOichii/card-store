package controller_test

import (
	"github.com/stretchr/testify/mock"
	"store.api/dto"
	"store.api/model"
	"store.api/query"
	"store.api/service"
)

type MockAuthService struct {
	mock.Mock
}

func newMockAuthService() *MockAuthService {
	return new(MockAuthService)
}

func (ser *MockAuthService) Login(user *dto.LoginDetails) (*dto.PrivateUserInfo, error) {
	args := ser.Called(user)
	switch user := args.Get(0).(type) {
	case *dto.PrivateUserInfo:
		return user, args.Error(1)
	case nil:
		return nil, args.Error(1)
	}
	return nil, args.Error(1)
}

func (ser *MockAuthService) Register(user *dto.RegisterDetails) error {
	args := ser.Called(user)
	return args.Error(0)
}

type MockCardService struct {
	mock.Mock
}

func newMockCardService() *MockCardService {
	return new(MockCardService)
}

func (ser *MockCardService) Add(c *dto.PostCard, posterId uint) (*dto.GetCard, error) {
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

func (ser *MockCardService) Query(query *query.CardQuery) *service.CardQueryResult {
	args := ser.Called(query)
	switch result := args.Get(0).(type) {
	case *service.CardQueryResult:
		return result
	case nil:
		return nil
	}
	return nil
}

func (ser *MockCardService) Update(c *dto.PostCard, cardId uint) (*dto.GetCard, error) {
	args := ser.Called(c, cardId)
	switch card := args.Get(0).(type) {
	case *dto.GetCard:
		return card, args.Error(1)
	case nil:
		return nil, args.Error(1)
	}
	return nil, args.Error(1)
}

func (ser *MockCardService) UpdatePrice(id uint, update *dto.PriceUpdate) (*dto.GetCard, error) {
	args := ser.Called(id, update)
	switch card := args.Get(0).(type) {
	case *dto.GetCard:
		return card, args.Error(1)
	case nil:
		return nil, args.Error(1)
	}
	return nil, args.Error(1)
}

func (ser *MockCardService) UpdateInStockAmount(id uint, update *dto.StockedAmountUpdate) (*dto.GetCard, error) {
	args := ser.Called(id, update)
	switch card := args.Get(0).(type) {
	case *dto.GetCard:
		return card, args.Error(1)
	case nil:
		return nil, args.Error(1)
	}
	return nil, args.Error(1)
}

func (ser *MockCardService) Languages() []*model.Language {
	args := ser.Called()
	return args.Get(0).([]*model.Language)
}

func (ser *MockCardService) Expansions() []*model.Expansion {
	args := ser.Called()
	return args.Get(0).([]*model.Expansion)
}

type MockCollectionService struct {
	mock.Mock
}

func newMockCollectionService() *MockCollectionService {
	return new(MockCollectionService)
}

func (ser *MockCollectionService) GetAll(userId uint) []*dto.GetCollection {
	args := ser.Called(userId)
	return args.Get(0).([]*dto.GetCollection)
}

func (ser *MockCollectionService) Create(c *dto.PostCollection, userId uint) (*dto.GetCollection, error) {
	args := ser.Called(c, userId)
	switch col := args.Get(0).(type) {
	case *dto.GetCollection:
		return col, args.Error(1)
	case nil:
		return nil, args.Error(1)
	}
	return nil, args.Error(1)
}

func (ser *MockCollectionService) EditSlot(cs *dto.PostCollectionSlot, colId uint, userId uint) (*dto.GetCollection, error) {
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

func (ser *MockCollectionService) Delete(id uint, userId uint) error {
	args := ser.Called(id, userId)
	return args.Error(0)
}

func (ser *MockCollectionService) UpdateInfo(col *dto.PostCollection, id uint, userId uint) (*dto.GetCollection, error) {
	args := ser.Called(col, id, userId)
	switch col := args.Get(0).(type) {
	case *dto.GetCollection:
		return col, args.Error(1)
	case nil:
		return nil, args.Error(1)
	}
	return nil, args.Error(1)
}

type MockCartService struct {
	mock.Mock
}

func newMockCartService() *MockCartService {
	return new(MockCartService)
}

func (ser *MockCartService) Get(userId uint) (*dto.GetCart, error) {
	args := ser.Called(userId)
	switch cart := args.Get(0).(type) {
	case *dto.GetCart:
		return cart, args.Error(1)
	case nil:
		return nil, args.Error(1)
	}
	return nil, args.Error(1)
}

func (ser *MockCartService) EditSlot(userId uint, cartSlot *dto.PostCartSlot) (*dto.GetCart, error) {
	args := ser.Called(userId, cartSlot)
	switch cart := args.Get(0).(type) {
	case *dto.GetCart:
		return cart, args.Error(1)
	case nil:
		return nil, args.Error(1)
	}
	return nil, args.Error(1)
}

type MockUserService struct {
	mock.Mock
}

func newMockUserService() *MockUserService {
	return new(MockUserService)
}

func (ser *MockUserService) ById(id uint) (*dto.PrivateUserInfo, error) {
	args := ser.Called(id)
	switch user := args.Get(0).(type) {
	case *dto.PrivateUserInfo:
		return user, args.Error(1)
	case nil:
		return nil, args.Error(1)
	}
	return nil, args.Error(1)
}

// ! duplicated from test/service/mocks_test.go
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
