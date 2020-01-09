package services

import (
	"../../entity"
	"../../users"
)

// CategoryService implements menu.CategoryService interface
type UserService struct {
	userRepo users.UserRepo
}

// NewUserService will create new UserService object
func NewUserService(UsrRepo users.UserRepo) *UserService {
	return &UserService{userRepo: UsrRepo}
}

// StoreUser persists new user information
func (us *UserService) StoreUser(user entity.User) error {

	err := us.userRepo.Registration(user)

	if err != nil {
		return err
	}

	return nil
}
