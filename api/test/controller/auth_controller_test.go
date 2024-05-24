package controller_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"store.api/auth"
	"store.api/config"
	"store.api/controller"
	"store.api/dto"
	"store.api/repository"
	"store.api/service"
)

func createAuthController(service service.UserService, repo repository.UserRepository) *controller.AuthController {
	return controller.NewAuthController(
		service,
		auth.NewJwtMiddleware(&config.Configuration{
			AuthKey: "test secret key",
		}, service, repo).Middle.LoginHandler,
	)
}

func Test_Auth_ShouldRegister(t *testing.T) {
	// arrange
	repo := createMockUserRepository()
	service := createMockUserService()
	controller := createAuthController(service, repo)
	data := dto.RegisterDetails{
		Username: "user",
		Password: "password",
		Email:    "mail@mail.com",
	}
	service.On("Register", mock.Anything).Return(nil)

	c, w := createTestContext(data)

	// act
	controller.Register(c)

	// assert
	assert.Equal(t, 200, w.Code)
}

func Test_Auth_ShouldNotRegister(t *testing.T) {
	// arrange
	repo := createMockUserRepository()
	service := createMockUserService()
	controller := createAuthController(service, repo)
	data := dto.RegisterDetails{
		Username: "user",
		Password: "password",
		Email:    "mail@mail.com",
	}
	service.On("Register", mock.Anything).Return(errors.New(""))

	c, w := createTestContext(data)

	// act
	controller.Register(c)

	// assert
	assert.Equal(t, 400, w.Code)
}

func Test_Auth_ShouldNotRegisterBadData(t *testing.T) {
	// arrange
	repo := createMockUserRepository()
	service := createMockUserService()
	controller := createAuthController(service, repo)
	data := []string{"first", "second"}
	c, w := createTestContext(data)

	// act
	controller.Register(c)

	// assert
	assert.Equal(t, 400, w.Code)
}

func Test_Auth_ShouldLogin(t *testing.T) {
	// arrange
	repo := createMockUserRepository()
	service := createMockUserService()
	controller := createAuthController(service, repo)
	data := dto.LoginDetails{
		Username: "user",
		Password: "password",
	}
	service.On("Login", mock.Anything).Return(&dto.PrivateUserInfo{
		Username: "user",
	}, nil)

	c, w := createTestContext(data)

	// act
	controller.Login(c)

	// assert
	assert.Equal(t, 200, w.Code)
}

func Test_Auth_ShouldNotLogin(t *testing.T) {
	// arrange
	repo := createMockUserRepository()
	service := createMockUserService()
	controller := createAuthController(service, repo)
	data := dto.LoginDetails{
		Username: "user",
		Password: "password",
	}
	service.On("Login", mock.Anything).Return(nil, errors.New(""))

	c, w := createTestContext(data)

	// act
	controller.Login(c)

	// assert
	assert.Equal(t, 401, w.Code)
}
