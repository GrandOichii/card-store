package controller

import "fmt"

func userNotFound(id uint) error {
	return fmt.Errorf("no user with id %d", id)
}
