package users

import "../entity"

type UserService interface {
	StoreUser(user entity.User) error
}
