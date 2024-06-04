package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"store.api/model"
	"store.api/service"
)

func newUserService(userRepo *MockUserRepository) service.UserService {
	return service.NewUserServiceImpl(
		userRepo,
	)
}

func Test_User_ShouldGetById(t *testing.T) {
	// arrange
	userRepo := newMockUserRepository()
	service := newUserService(userRepo)

	userRepo.On("FindById", mock.Anything).Return(&model.User{})

	// act
	user, err := service.ById(1)

	// assert
	assert.NotNil(t, user)
	assert.Nil(t, err)
}

func Test_User_ShouldNotGetById(t *testing.T) {
	// arrange
	userRepo := newMockUserRepository()
	s := newUserService(userRepo)

	userRepo.On("FindById", mock.Anything).Return(nil)

	// act
	user, err := s.ById(1)

	// assert
	assert.Nil(t, user)
	assert.Equal(t, service.ErrUserNotFound, err)
}
