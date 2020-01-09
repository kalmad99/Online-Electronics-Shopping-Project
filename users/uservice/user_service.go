package uservice

import (
	"../../entity"
	"../../users"
)

// CategoryService implements menu.CategoryService interface
type UserService struct {
	userRepo users.UserRepository
}

// NewUserService will create new UserService object
func NewUserService(UsrRepo users.UserRepository) *UserService {
	return &UserService{userRepo: UsrRepo}
}

// Users returns list of users
func (us *UserService) Users() ([]entity.User, error) {

	users, err := us.userRepo.Users()

	if err != nil {
		return nil, err
	}

	return users, nil
}

// StoreUser persists new user information
func (us *UserService) StoreUser(user entity.User) error {

	err := us.userRepo.StoreUser(user)

	if err != nil {
		return err
	}

	return nil
}

// User returns a user object with a given id
func (us *UserService) User(email string) (entity.User, error) {

	u, err := us.userRepo.Login(email)

	if err != nil {
		return u, err
	}

	return u, nil
}

func (us *UserService) UserwithID(id int) (entity.User, error) {

	u, err := us.userRepo.UserwithID(id)

	if err != nil {
		return u, err
	}

	return u, nil
}

// UpdateUser updates a user with new data
func (us *UserService) UpdateUser(user entity.User) error {

	err := us.userRepo.UpdateUser(user)

	if err != nil {
		return err
	}

	return nil
}

// DeleteUser delete a user by its id
func (us *UserService) DeleteUser(id int) error {

	err := us.userRepo.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}
