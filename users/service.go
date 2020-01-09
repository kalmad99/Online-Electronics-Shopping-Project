package users

import "../entity"

type UserService interface {
	Users() ([]entity.User, error)
	Login(email string) (entity.User, error)
	UserwithID(id int) (entity.User, error)
	UpdateUser(user entity.User) error
	DeleteUser(id int) error
	StoreUser(user entity.User) error
}
