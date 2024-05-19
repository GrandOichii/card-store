package service

import "store.api/dto"

type UserService interface {
	Register(*dto.RegisterDetails) error
	Login(*dto.LoginDetails) (*dto.PrivateUserInfo, error)
}
