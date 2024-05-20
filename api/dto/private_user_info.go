package dto

import "store.api/model"

type PrivateUserInfo struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
}

func NewPrivateUserInfo(user *model.User) *PrivateUserInfo {
	result := PrivateUserInfo{
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	}

	return &result
}
