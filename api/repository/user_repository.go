package repository

import "store.api/model"

type UserRepository interface {
	Save(*model.User) error
	FindByUsername(username string) *model.User
	FindByEmail(email string) *model.User
}
