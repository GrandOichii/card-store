package dto

import "store.api/model"

type PrivateUserInfo struct {
	Username string `json:"username"`
}

func NewPrivateUserInfo(user *model.User) *PrivateUserInfo {
	result := PrivateUserInfo{
		Username: user.Username,
	}

	return &result
}
