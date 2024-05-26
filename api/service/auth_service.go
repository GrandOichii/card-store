package service

import "store.api/dto"

type AuthService interface {
	Register(*dto.RegisterDetails) error
	Login(*dto.LoginDetails) (*dto.PrivateUserInfo, error)
}
