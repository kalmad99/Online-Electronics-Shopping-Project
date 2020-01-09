package services

import (
	"../../entity"
	"Online-Electronics-Shopping-Project/users"
)

// CategoryService implements menu.CategoryService interface
type UserService struct {
	userRepo users.UserRepository
}

// NewUserService will create new UserService object
func NewUserService(UsrRepo users.UserRepository) *UserService {
	return &UserService{userRepo: UsrRepo}
}

// StoreUser persists new user information
func (us *UserService) StoreUser(user entity.User) error {

	err := us.userRepo.StoreUser(user)

	if err != nil {
		return err
	}

	return nil
}
