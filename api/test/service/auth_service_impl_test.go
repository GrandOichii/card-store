package service_test

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"store.api/dto"
	"store.api/model"
	"store.api/security"
	"store.api/service"
)

func createAuthService(userRepo *MockUserRepository, cartRepo *MockCartRepository) service.AuthService {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return service.NewAuthServiceImpl(
		userRepo,
		cartRepo,
		validate,
	)
}

func Test_User_ShouldRegister(t *testing.T) {
	// arrange
	userRepo := newMockUserRepository()
	cartRepo := newMockTaskRepository()
	service := createAuthService(userRepo, cartRepo)
	data := dto.RegisterDetails{
		Username: "user",
		Password: "password",
		Email:    "mail@mail.com",
	}

	userRepo.On("Save", mock.Anything).Return(nil)
	userRepo.On("FindByUsername", data.Username).Return(nil)
	userRepo.On("FindByEmail", data.Email).Return(nil)
	cartRepo.On("Save", mock.Anything).Return(nil)

	// act
	err := service.Register(&data)

	// assert
	assert.Nil(t, err)
}

func Test_User_ShouldNotRegisterUsernameTaken(t *testing.T) {
	// arrange
	userRepo := newMockUserRepository()
	cartRepo := newMockTaskRepository()
	service := createAuthService(userRepo, cartRepo)
	data := dto.RegisterDetails{
		Username: "user",
		Password: "password",
		Email:    "mail@mail.com",
	}

	userRepo.On("FindByUsername", data.Username).Return(&model.User{})
	userRepo.On("FindByEmail", data.Email).Return(nil)

	// act
	err := service.Register(&data)

	// assert
	assert.NotNil(t, err)
}

func Test_User_ShouldNotRegisterEmailTaken(t *testing.T) {
	// arrange
	userRepo := newMockUserRepository()
	cartRepo := newMockTaskRepository()
	service := createAuthService(userRepo, cartRepo)
	email := "mail@mail.com"
	data := dto.RegisterDetails{
		Username: "user",
		Password: "password",
		Email:    email,
	}

	userRepo.On("FindByUsername", data.Username).Return(nil)
	userRepo.On("FindByEmail", data.Email).Return(&model.User{
		Email:    "mail@mail.com",
		Verified: true,
	})

	// act
	err := service.Register(&data)

	// assert
	assert.NotNil(t, err)
}

func Test_User_ShouldNotLogin(t *testing.T) {
	// arrange
	userRepo := newMockUserRepository()
	cartRepo := newMockTaskRepository()
	service := createAuthService(userRepo, cartRepo)
	data := dto.LoginDetails{
		Username: "user",
		Password: "password",
	}

	userRepo.On("FindByUsername", data.Username).Return(nil)

	// act
	login, err := service.Login(&data)

	// assert
	assert.Nil(t, login)
	assert.NotNil(t, err)
}

func Test_User_ShouldNotLoginIncorrectPassword(t *testing.T) {
	// arrange
	userRepo := newMockUserRepository()
	cartRepo := newMockTaskRepository()
	service := createAuthService(userRepo, cartRepo)
	data := dto.LoginDetails{
		Username: "user",
		Password: "password",
	}
	existing := model.User{
		Username:     data.Username,
		PasswordHash: "passwordHash",
	}

	userRepo.On("FindByUsername", data.Username).Return(&existing)

	// act
	login, err := service.Login(&data)

	// assert
	assert.Nil(t, login)
	assert.NotNil(t, err)
}

func Test_User_ShouldLogin(t *testing.T) {
	// arrange
	userRepo := newMockUserRepository()
	cartRepo := newMockTaskRepository()
	service := createAuthService(userRepo, cartRepo)
	data := dto.LoginDetails{
		Username: "user",
		Password: "password",
	}
	hash, _ := security.HashPassword(data.Password)
	existing := model.User{
		Username:     data.Username,
		PasswordHash: hash,
	}

	userRepo.On("FindByUsername", data.Username).Return(&existing)

	// act
	login, err := service.Login(&data)

	// assert
	assert.NotNil(t, login)
	assert.Nil(t, err)
}
